package main

import (
	"encoding/xml"
	"io/ioutil"
)

const (
	TagTestSuites = "testsuites"
	TagTestSuite  = "testsuite"
	TagTestCase   = "testcase"
	TagFailure    = "failure"
	TagSkipped    = "skipped"
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

type JUnitSuit struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Tests     int             `xml:"tests,attr"`
	Errors    int             `xml:"errors,attr"`
	Failures  int             `xml:"failures,attr"`
	Skip      int             `xml:"skip,attr"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

type JUnitTestResult struct {
	XMLName xml.Name    `xml:"testsuites"`
	Suites  []JUnitSuit `xml:"testsuite"`
}

func writeToXML(result JUnitTestResult, filename string) error {
	output, err := xml.MarshalIndent(result, "  ", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, append([]byte(xml.Header), output...), 0755)
}

func xmlName(name string) xml.Name {
	return xml.Name{Local: name}
}
