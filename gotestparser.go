package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type ParserState int8

const (
	parserNextStep = iota
	parserTestCaseResult
	parserTestCaseMessage
	parserTestSuitSummary
	parserFail
)

const (
	msgFailureType    = "gotest.error"
	msgFailureMessage = "error"
	msgSkippedType    = "gotest.skipped"
	msgSkippedMessage = "skipped"
)

var errParserError = func(err string) error { return fmt.Errorf("Error while parsing at line : %s", err) }

var (
	rgxTestCase         = regexp.MustCompile(`^===\sRUN\s(?:[A-Za-z0-9_-]+)\z`)
	rgxResultCase       = regexp.MustCompile(`^---\s(PASS|FAIL|SKIP):\s([A-Za-z0-9_-]+)\s\((\d+\.\d{2})\sseconds\)\z`)
	rgxSuiteResult      = regexp.MustCompile(`^(PASS|FAIL)\z`)
	rgxTestSuiteSummary = regexp.MustCompile(`^(?:ok|FAIL|\?)\s+([A-Za-z0-9_\-/\\\.]+)\s+(?:(?:\d+\.\d+s)|(?:\[no\stest\sfiles\]))\z`)
	rgxEmptySuite       = regexp.MustCompile(`^?\s+([A-Za-z0-9_\-/\\\.]+)\s+\[no\stest\sfiles\]\z`)
	rgxExitStatus       = regexp.MustCompile(`^exit\s+status\s+\d+\z`)
)

/*
  (1) Next step (default)
      a) test case : "=== RUN TestPackage_1" -> (2)
      b) test suit summary : "PASS" or "FAIL" -> (4)
      c) empty test suit : "?   gojunitxml/package_6	[no test files]" -> (4c)

  (2) Parsing testcase result
      a) test case summary : "--- PASS: TestPackage_1 (0.00 seconds)" -> (1)
                             "--- FAIL: TestPackage_2 (0.00 seconds)" -> (3a)
                             "--- SKIP: TestPackage_2 (0.00 seconds)" -> (3a)

  (3) Parsing message - SKIP or FAIL have a message (could be multi-line)
      a) error message    : "package_1_test.go:13: Some error message" -> (3)
      b) test case        : "=== RUN TestPackage_1" -> (2)
      c) testsuit summary : "FAIL" -> (4b)

  (4) Parsing testsuite summary
      a) success : "ok  	gojunitxml/package_2	0.006s" -> (1)
      b) fail    : "FAIL	gojunitxml/package_1	0.005s" -> (1)
      c) skipped : "?   	gojunitxml/package_6	[no test files]" -> (1)


    ** "exit status" line is ignored, "build failed" is not.
*/

func parseGoTest(reader io.Reader) (*JUnitTestResult, error) {
	var (
		testResult = newJUnitTestResult(tagTestSuites)
		testSuite  = newJUnitSuite(tagTestSuite)
		state      = ParserState(parserNextStep)
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if state = parseLine(scanner.Text(), state, testSuite, testResult); state == parserFail {
			return nil, errParserError(scanner.Text())
		}
	}
	return testResult, nil
}

func parseLine(line string, state ParserState, testSuite *JUnitSuite, testResult *JUnitTestResult) ParserState {
	switch state {
	case parserNextStep:
		return parseNextStep(line, testSuite, testResult)
	case parserTestCaseResult:
		return parseTestCaseResult(line, testSuite)
	case parserTestCaseMessage:
		return parseTestCaseMessage(line, testSuite)
	case parserTestSuitSummary:
		return parseTestSuiteSummary(line, testSuite, testResult)
	default:
		return parserFail
	}
}

func parseNextStep(line string, testSuite *JUnitSuite, testResult *JUnitTestResult) ParserState {
	switch {
	case rgxTestCase.MatchString(line):
		return parserTestCaseResult
	case rgxSuiteResult.MatchString(line):
		return parserTestSuitSummary
	case rgxEmptySuite.MatchString(line):
		return parseTestSuiteSummary(line, testSuite, testResult)
	default:
		return parserFail
	}
}

func parseTestCaseResult(line string, testSuite *JUnitSuite) ParserState {
	var (
		match    []string
		testCase *JUnitTestCase
		state    ParserState
	)

	if match = rgxResultCase.FindStringSubmatch(line); match == nil {
		return parserFail
	}

	switch match[1] {
	case "FAIL":
		state = parserTestCaseMessage
		testSuite.Failures += 1
		testCase = newJUnitTestCase(
			tagTestCase, match[2], match[3],
			newJUnitTestCaseMessage(tagFailure, msgFailureMessage, msgFailureType))
	case "SKIP":
		state = parserTestCaseMessage
		testSuite.Skip += 1
		testCase = newJUnitTestCase(
			tagTestCase, match[2], match[3],
			newJUnitTestCaseMessage(tagSkipped, msgSkippedMessage, msgSkippedType))
	default:
		state = parserNextStep
		testCase = newJUnitTestCase(
			tagTestCase, match[2], match[3], nil)
	}

	testSuite.Tests += 1
	testSuite.TestCases = append(testSuite.TestCases, *testCase)
	return state
}

func parseTestCaseMessage(line string, testSuite *JUnitSuite) ParserState {
	if rgxTestCase.MatchString(line) {
		return parserTestCaseResult
	}
	if rgxSuiteResult.MatchString(line) {
		return parserTestSuitSummary
	}
	testSuite.TestCases[len(testSuite.TestCases)-1].Message.Content += strings.Trim(line, " \t") + ";"
	return parserTestCaseMessage
}

func parseTestSuiteSummary(line string, testSuite *JUnitSuite, testResult *JUnitTestResult) ParserState {
	var match []string

	if rgxExitStatus.MatchString(line) {
		return parserTestSuitSummary
	}
	if match = rgxTestSuiteSummary.FindStringSubmatch(line); match == nil {
		return parserFail
	}

	testSuite.Name = match[1]
	for i := range testSuite.TestCases {
		testSuite.TestCases[i].ClassName = testSuite.Name
	}
	testResult.Suites = append(testResult.Suites, *testSuite)
	*testSuite = *newJUnitSuite(tagTestSuite)
	return parserNextStep
}
