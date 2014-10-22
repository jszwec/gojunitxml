package gojunitxml

import "encoding/xml"

type result uint8

const (
	passed result = iota
	failed
	skipped
	unknown
)

const (
	msgFailure        = "failure"
	msgFailureType    = "gotest.error"
	msgFailureMessage = "error"
	msgSkipped        = "skipped"
	msgSkippedType    = "gotest.skipped"
	msgSkippedMessage = "skipped"
)

var resmap = map[string]result{
	"PASS": passed,
	"FAIL": failed,
	"SKIP": skipped,
}

type message struct {
	XMLName xml.Name
	Type    string `xml:"type,attr"`
	Message string `xml:"message,attr"`
	Content string `xml:",chardata"`
}

type testcase struct {
	ClassName string     `xml:"classname,attr"`
	Name      string     `xml:"name,attr"`
	Time      string     `xml:"time,attr"`
	Messages  []*message `xml:",any"`
	result    result
}

type testsuite struct {
	Name      string     `xml:"name,attr"`
	Tests     int        `xml:"tests,attr"`
	Errors    int        `xml:"errors,attr"`
	Failures  int        `xml:"failures,attr"`
	Skip      int        `xml:"skip,attr"`
	TestCases []testcase `xml:"testcase"`
}

type testsuites struct {
	Suites []testsuite `xml:"testsuite"`
}

func (t testsuites) Marshal() (b []byte, err error) {
	if b, err = xml.MarshalIndent(t, "  ", "    "); err == nil {
		b = append([]byte(xml.Header), b...)
	}
	return
}

func newMessageResult(cnt string, res result) (m *message) {
	switch res {
	case failed:
		m = newMessage(msgFailure, msgFailureMessage, msgFailureType, cnt)
	case skipped:
		m = newMessage(msgSkipped, msgSkippedMessage, msgSkippedType, cnt)
	}
	return
}

func newMessage(tag, msg, typ, cnt string) *message {
	return &message{
		XMLName: xml.Name{Local: tag},
		Message: msg,
		Type:    typ,
		Content: cnt,
	}
}

func resultString(s string) result {
	if r, ok := resmap[s]; ok {
		return r
	}
	return unknown
}
