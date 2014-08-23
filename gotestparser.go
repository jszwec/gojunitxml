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
	ParserNextStep = iota
	ParserTestCaseResult
	ParserTestCaseMessage
	ParserTestSuitSummary
	ParserFail
)

const (
	MsgFailureType    = "gotest.error"
	MsgFailureMessage = "error"
	MsgSkippedType    = "gotest.skipped"
	MsgSkippedMessage = "skipped"
)

var ErrParserError = func(err string) error { return fmt.Errorf("Error while parsing at line : %s", err) }

var (
	RgxTestCase         = regexp.MustCompile(`^===\sRUN\s(?:[A-Za-z0-9_-]+)\z`)
	RgxResultCase       = regexp.MustCompile(`^---\s(PASS|FAIL|SKIP):\s([A-Za-z0-9_-]+)\s\((\d+\.\d{2})\sseconds\)\z`)
	RgxSuiteResult      = regexp.MustCompile(`^(PASS|FAIL)\z`)
	RgxTestSuiteSummary = regexp.MustCompile(`^(?:ok|FAIL|\?)\s+([A-Za-z0-9_\-/\\\.]+)\s+(?:(?:\d+\.\d+s)|(?:\[no\stest\sfiles\]))\z`)
	RgxEmptySuite       = regexp.MustCompile(`^?\s+([A-Za-z0-9_\-/\\\.]+)\s+\[no\stest\sfiles\]\z`)
	RgxExitStatus       = regexp.MustCompile(`^exit\s+status\s+\d+\z`)
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

func parseGoTest(reader io.Reader) (JUnitTestResult, error) {
	var (
		testResult JUnitTestResult = JUnitTestResult{XMLName: xmlName(TagTestSuites)}
		testSuit   JUnitSuit       = JUnitSuit{XMLName: xmlName(TagTestSuite)}
		state      ParserState     = ParserNextStep
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if state = parseLine(scanner.Text(), state, &testSuit, &testResult); state == ParserFail {
			return JUnitTestResult{}, ErrParserError(scanner.Text())
		}
	}
	return testResult, nil
}

func parseLine(line string, state ParserState, testSuit *JUnitSuit, testResult *JUnitTestResult) ParserState {
	switch state {
	case ParserNextStep:
		return parseNextStep(line, testSuit, testResult)
	case ParserTestCaseResult:
		return parseTestCaseResult(line, testSuit)
	case ParserTestCaseMessage:
		return parseTestCaseMessage(line, testSuit)
	case ParserTestSuitSummary:
		return parseTestSuiteSummary(line, testSuit, testResult)
	default:
		return ParserFail
	}
}

func parseNextStep(line string, testSuit *JUnitSuit, testResult *JUnitTestResult) ParserState {
	switch {
	case RgxTestCase.MatchString(line):
		return ParserTestCaseResult
	case RgxSuiteResult.MatchString(line):
		return ParserTestSuitSummary
	case RgxEmptySuite.MatchString(line):
		return parseTestSuiteSummary(line, testSuit, testResult)
	default:
		return ParserFail
	}
}

func parseTestCaseResult(line string, testSuit *JUnitSuit) ParserState {
	var (
		match    []string       = RgxResultCase.FindStringSubmatch(line)
		testCase *JUnitTestCase = &JUnitTestCase{XMLName: xmlName(TagTestCase)}
		state    ParserState    = ParserNextStep
	)

	if match == nil {
		return ParserFail
	}

	testCase.Name = match[2]
	testCase.Time = match[3]
	testSuit.Tests += 1

	switch match[1] {
	case "FAIL":
		testSuit.Failures += 1
		testCase.Message = &JUnitTestCaseMessage{XMLName: xmlName(TagFailure), Message: MsgFailureMessage, Type: MsgFailureType}
		state = ParserTestCaseMessage
	case "SKIP":
		testSuit.Skip += 1
		testCase.Message = &JUnitTestCaseMessage{XMLName: xmlName(TagSkipped), Message: MsgSkippedMessage, Type: MsgSkippedType}
		state = ParserTestCaseMessage
	default:
		break
	}

	testSuit.TestCases = append(testSuit.TestCases, *testCase)
	return state
}

func parseTestCaseMessage(line string, testSuit *JUnitSuit) ParserState {
	if RgxTestCase.MatchString(line) {
		return ParserTestCaseResult
	}
	if RgxSuiteResult.MatchString(line) {
		return ParserTestSuitSummary
	}

	testSuit.TestCases[len(testSuit.TestCases)-1].Message.Content += strings.Trim(line, " \t") + ";"
	return ParserTestCaseMessage
}

func parseTestSuiteSummary(line string, testSuit *JUnitSuit, testResult *JUnitTestResult) ParserState {
	if RgxExitStatus.MatchString(line) {
		return ParserTestSuitSummary
	}
	var match []string = RgxTestSuiteSummary.FindStringSubmatch(line)

	if match == nil {
		return ParserFail
	}

	testSuit.Name = match[1]
	testResult.Suites = append(testResult.Suites, *testSuit)
	for i := range testSuit.TestCases {
		testSuit.TestCases[i].ClassName = match[1]
	}
	*testSuit = JUnitSuit{XMLName: xmlName(TagTestSuite)}
	return ParserNextStep
}
