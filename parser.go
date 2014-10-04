package gojunitxml

import (
	"bufio"
	"io"
	"regexp"
	"strings"
)

const (
	msgFailureType    = "gotest.error"
	msgFailureMessage = "error"
	msgSkippedType    = "gotest.skipped"
	msgSkippedMessage = "skipped"
)

var (
	reTestCase    = regexp.MustCompile(`^===\sRUN\s([A-Za-z0-9_-]+)\z`)
	reResultCase  = regexp.MustCompile(`^---\s(PASS|FAIL|SKIP):\s([A-Za-z0-9_-]+)\s\((\d+\.\d{2})\sseconds\)\z`)
	reResultSuite = regexp.MustCompile(`^(PASS|FAIL)\z`)
	reTestSuite   = regexp.MustCompile(`^(?:ok|FAIL|\?)\s+([A-Za-z0-9_\-/\\\.]+)\s+(?:(?:\d+\.\d+s)|(?:\[no\stest\sfiles\]))\z`)
	reExitStatus  = regexp.MustCompile(`^exit\s+status\s+\d+\z`)
)

type parser struct {
	ts  *testsuites
	tcs []testcase
}

func newparser() *parser {
	return &parser{
		ts: &testsuites{},
	}
}

func result(s string) caseresult {
	switch s {
	case "PASS":
		return passed
	case "FAIL":
		return failed
	case "SKIP":
		return skipped
	default:
		return unknown
	}
}

func (p *parser) addtestcase(m []string) {
	p.tcs = append(p.tcs, testcase{
		Name:   m[2],
		Time:   m[3],
		result: result(m[1])})
}

func (p *parser) addmessage(t *testcase, line string) {
	switch t.result {
	case failed:
		t.Messages = append(t.Messages, errorMessage(line))
	case skipped:
		t.Messages = append(t.Messages, skipMessage(line))
	default:
		break
	}
}

func (p *parser) addtestsuite(m []string) {
	suite := testsuite{
		Name: func(s string) (n string) {
			n = strings.Replace(m[1], ".", "_", -1)
			i := strings.LastIndex(m[1], "/")
			if i > -1 && i < len(n)-1 {
				n += "." + n[i+1:]
			} else {
				n += "." + n
			}
			return
		}(m[1]),
		Tests: len(p.tcs),
	}
	for _, tc := range p.tcs {
		tc.ClassName = suite.Name
		switch tc.result {
		case failed:
			suite.Failures++
		case skipped:
			suite.Skip++
		default:
			break
		}
		suite.TestCases = append(suite.TestCases, tc)
	}
	p.ts.Suites = append(p.ts.Suites, suite)
	p.tcs = p.tcs[:0]
}

func (p *parser) parseLine(line string) {
	switch {
	case reTestCase.MatchString(line), reResultSuite.MatchString(line),
		reExitStatus.MatchString(line):
	case reResultCase.MatchString(line):
		p.addtestcase(reResultCase.FindStringSubmatch(line))
	case reTestSuite.MatchString(line):
		p.addtestsuite(reTestSuite.FindStringSubmatch(line))
	case len(p.tcs) > 0:
		p.addmessage(&p.tcs[len(p.tcs)-1], line)
	}
}

// Parse : Parses given "go test -v" output.
// Parser will ignore any lines which are not related to gotest.
func Parse(reader io.Reader) *testsuites {
	p, scanner := newparser(), bufio.NewScanner(reader)
	for scanner.Scan() {
		p.parseLine(strings.Trim(scanner.Text(), "\t "))
	}
	return p.ts
}
