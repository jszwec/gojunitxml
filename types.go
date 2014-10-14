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

func (t testsuites) Marshal() ([]byte, error) {
	b, err := xml.MarshalIndent(t, "  ", "    ")
	if err != nil {
		return nil, err
	}
	return append([]byte(xml.Header), b...), nil
}

func newMessageResult(cnt string, res result) *message {
	switch res {
	case failed:
		return newMessage(msgFailure, msgFailureMessage, msgFailureType, cnt)
	case skipped:
		return newMessage(msgSkipped, msgSkippedMessage, msgSkippedType, cnt)
	default:
		return nil
	}
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
