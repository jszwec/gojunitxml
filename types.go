package main

import (
	"encoding/xml"
	"io/ioutil"
)

const (
	tagTestSuites = "testsuites"
	tagTestSuite  = "testsuite"
	tagTestCase   = "testcase"
	tagFailure    = "failure"
	tagSkipped    = "skipped"
)

type JUnitTestCaseMessage struct {
	XMLName xml.Name
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Content string `xml:",chardata"`
}

type JUnitTestCase struct {
	XMLName   xml.Name              `xml:"testcase"`
	ClassName string                `xml:"classname,attr"`
	Name      string                `xml:"name,attr"`
	Time      string                `xml:"time,attr"`
	Message   *JUnitTestCaseMessage `xml:",any"`
}

type JUnitSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Errors    int             `xml:"errors,attr"`
	Failures  int             `xml:"failures,attr"`
	Skip      int             `xml:"skip,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

type JUnitTestResult struct {
	XMLName xml.Name     `xml:"testsuites"`
	Suites  []JUnitSuite `xml:"testsuite"`
}

func newJUnitTestResult(tag string) *JUnitTestResult {
	return &JUnitTestResult{XMLName: xml.Name{Local: tag}}
}

func newJUnitSuite(tag string) *JUnitSuite {
	return &JUnitSuite{XMLName: xml.Name{Local: tag}}
}

func newJUnitTestCaseMessage(tag, message, typ string) *JUnitTestCaseMessage {
	return &JUnitTestCaseMessage{
		XMLName: xml.Name{Local: tag},
		Message: message,
		Type:    typ,
	}
}

func newJUnitTestCase(tag, name, time string, message *JUnitTestCaseMessage) *JUnitTestCase {
	return &JUnitTestCase{
		XMLName: xml.Name{Local: tag},
		Name:    name,
		Time:    time,
		Message: message,
	}
}

func (r *JUnitTestResult) WriteToXML(filename string) error {
	output, err := xml.MarshalIndent(r, "  ", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, append([]byte(xml.Header), output...), 0755)
}
