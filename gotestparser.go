package main

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type PARSING_STATE int8

const (
	PARSING_NEXT_STEP = iota
	PARSING_TEST_CASE_RESULT
	PARSING_TEST_CASE_MESSAGE
	PARSING_TEST_SUIT_SUMMARY
	FAIL
)

const (
	MSG_PARSING_ERROR   = "Error while parsing at line : %s"
	MSG_FAILURE_TYPE    = "gotest.error"
	MSG_FAILURE_MESSAGE = "error"
	MSG_SKIPPED_TYPE    = "gotest.skipped"
	MSG_SKIPPED_MESSAGE = "skipped"
)

var (
	testCaseRegex         = regexp.MustCompile(`^===\sRUN\s(?:[A-Za-z0-9_-]+)\z`)
	resultCaseRegex       = regexp.MustCompile(`^---\s(PASS|FAIL|SKIP):\s([A-Za-z0-9_-]+)\s\((\d+\.\d{2})\sseconds\)\z`)
	testSuitResultRegex   = regexp.MustCompile(`^(PASS|FAIL)\z`)
	testSuiteSummaryRegex = regexp.MustCompile(`^(?:ok|FAIL|\?)\s+([A-Za-z0-9_\-/\\\.]+)\s+(?:(?:\d+\.\d+s)|(?:\[no\stest\sfiles\]))\z`)
	emptySuitRegex        = regexp.MustCompile(`^?\s+([A-Za-z0-9_\-/\\\.]+)\s+\[no\stest\sfiles\]\z`)
	exitstatusRegex       = regexp.MustCompile(`^exit\s+status\s+\d+\z`)
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
		testResult JUnitTestResult = JUnitTestResult{XMLName: xmlName(TAG_TESTSUITES)}
		testSuit   JUnitSuit       = JUnitSuit{XMLName: xmlName(TAG_TESTSUITE)}
		state      PARSING_STATE   = PARSING_NEXT_STEP
	)

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		if state = parseLine(scanner.Text(), state, &testSuit, &testResult); state == FAIL {
			return JUnitTestResult{}, fmt.Errorf(MSG_PARSING_ERROR, scanner.Text())
		}
	}
	return testResult, nil
}

func parseLine(line string, state PARSING_STATE, testSuit *JUnitSuit, testResult *JUnitTestResult) PARSING_STATE {
	switch state {
	case PARSING_NEXT_STEP:
		return parseNextStep(line, testSuit, testResult)
	case PARSING_TEST_CASE_RESULT:
		return parseTestCaseResult(line, testSuit)
	case PARSING_TEST_CASE_MESSAGE:
		return parseTestCaseMessage(line, testSuit)
	case PARSING_TEST_SUIT_SUMMARY:
		return parseTestSuiteSummary(line, testSuit, testResult)
	default:
		return FAIL
	}
}

func parseNextStep(line string, testSuit *JUnitSuit, testResult *JUnitTestResult) PARSING_STATE {
	if testCaseRegex.MatchString(line) {
		return PARSING_TEST_CASE_RESULT
	} else if testSuitResultRegex.MatchString(line) {
		return PARSING_TEST_SUIT_SUMMARY
	} else if emptySuitRegex.MatchString(line) {
		return parseTestSuiteSummary(line, testSuit, testResult)
	}
	return FAIL
}

func parseTestCaseResult(line string, testSuit *JUnitSuit) PARSING_STATE {
	var (
		match    []string       = resultCaseRegex.FindStringSubmatch(line)
		testCase *JUnitTestCase = &JUnitTestCase{XMLName: xmlName(TAG_TESTCASE)}
		state    PARSING_STATE  = PARSING_NEXT_STEP
	)

	if match == nil {
		return FAIL
	}

	testCase.Name = match[2]
	testCase.Time = match[3]
	testSuit.Tests += 1

	switch match[1] {
	case "FAIL":
		testSuit.Failures += 1
		testCase.Message = &JUnitTestCaseMessage{XMLName: xmlName(TAG_FAILURE), Message: MSG_FAILURE_MESSAGE, Type: MSG_FAILURE_TYPE}
		state = PARSING_TEST_CASE_MESSAGE
	case "SKIP":
		testSuit.Skip += 1
		testCase.Message = &JUnitTestCaseMessage{XMLName: xmlName(TAG_SKIPPED), Message: MSG_SKIPPED_MESSAGE, Type: MSG_SKIPPED_TYPE}
		state = PARSING_TEST_CASE_MESSAGE
	default:
		break
	}

	testSuit.TestCases = append(testSuit.TestCases, *testCase)
	return state
}

func parseTestCaseMessage(line string, testSuit *JUnitSuit) PARSING_STATE {
	if testCaseRegex.MatchString(line) {
		return PARSING_TEST_CASE_RESULT
	}
	if testSuitResultRegex.MatchString(line) {
		return PARSING_TEST_SUIT_SUMMARY
	}

	testSuit.TestCases[len(testSuit.TestCases)-1].Message.Content += strings.Trim(line, " \t") + ";"

	return PARSING_TEST_CASE_MESSAGE
}

func parseTestSuiteSummary(line string, testSuit *JUnitSuit, testResult *JUnitTestResult) PARSING_STATE {
	if exitstatusRegex.MatchString(line) {
		return PARSING_TEST_SUIT_SUMMARY
	}
	var match []string = testSuiteSummaryRegex.FindStringSubmatch(line)

	if match == nil {
		return FAIL
	}

	testSuit.Name = match[1]
	testResult.Suites = append(testResult.Suites, *testSuit)
	for i := range testSuit.TestCases {
		testSuit.TestCases[i].ClassName = match[1]
	}
	*testSuit = JUnitSuit{XMLName: xmlName(TAG_TESTSUITE)}
	return PARSING_NEXT_STEP
}
