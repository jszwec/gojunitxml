package gojunitxml

import "encoding/xml"

type caseresult uint8

const (
	passed caseresult = iota
	failed
	skipped
	unknown
)

const (
	msgFailureType    = "gotest.error"
	msgFailureMessage = "error"
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
	ClassName string    `xml:"classname,attr"`
	Name      string    `xml:"name,attr"`
	Time      string    `xml:"time,attr"`
	Messages  []message `xml:",any"`
	result    caseresult
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

func errorMessage(content string) message {
	return message{
		XMLName: xml.Name{Local: "failure"},
		Message: msgFailureMessage,
		Type:    msgFailureType,
		Content: content,
	}
}

func skipMessage(content string) message {
	return message{
		XMLName: xml.Name{Local: "skipped"},
		Message: msgSkippedMessage,
		Type:    msgSkippedType,
		Content: content,
	}
}
