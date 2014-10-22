package gojunitxml

import (
	"bufio"
	"io"
	"regexp"
	"strings"
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

func (p *parser) addtestcase(m []string) {
	p.tcs = append(p.tcs, testcase{
		Name:   m[2],
		Time:   m[3],
		result: resultString(m[1]),
	})
}

func (p *parser) addmessage(t *testcase, line string) {
	if m := newMessageResult(line, t.result); m != nil {
		t.Messages = append(t.Messages, m)
	}
}

func (p *parser) addtestsuite(m []string) {
	suite := testsuite{
		// testsuite's name will be slightly changed e.g.
		// github.com/jszwec/gojunitxml -> github_com/jszwec/gojunitxml.gojunitxml
		//
		// This is for Jenkins UI
		//   Package = github_com/jszwec/gojunitxml
		//   Class   = gojunitxml
		Name: func(s string) string {
			n := strings.Replace(m[1], ".", "_", -1)
			if i := strings.LastIndex(m[1], "/"); i > -1 && i < len(n)-1 {
				return n + "." + n[i+1:]
			}
			return n + "." + n
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
