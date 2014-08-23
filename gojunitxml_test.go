package main

import (
	"bufio"
	"encoding/xml"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func gotestDir(filename string) string {
	return filepath.Join("data", "gotest", filename)
}

func dataDir(filename string) string {
	return filepath.Join("data", filename)
}

func openTestFile(testfile string, t *testing.T) *os.File {
	file, err := os.Open(testfile)
	if err != nil {
		t.Fatal(err)
	}
	return file
}

func trim(val string) string {
	return strings.Trim(val, " \t\n")
}

func checkTestResult(data []byte, result JUnitTestResult, t *testing.T) {
	v := JUnitTestResult{Suites: nil}
	err := xml.Unmarshal(data, &v)
	if err != nil {
		t.Fatal("expected unsmarshalling to succeed")
	}
	if !reflect.DeepEqual(v, result) {
		t.Errorf("expected %v\nactual %v", v, result)
	}
}

func checkXMLOutput(outputfile string, data []string, t *testing.T) {
	var (
		testfile *os.File = openTestFile(outputfile, t)
		datasize int      = len(data)
		counter  int
	)
	defer testfile.Close()

	scanner := bufio.NewScanner(testfile)
	for counter = 0; scanner.Scan(); counter++ {
		if counter >= datasize {
			t.Fatal("XML output is bigger than expected data")
		}
		expected := trim(data[counter])
		actual := trim(scanner.Text())
		if expected != actual {
			t.Errorf("expected %s, actual %s", expected, actual)
		}
	}
	if counter != datasize {
		t.Error("expected datasize is smaller than XML output")
	}
}

func Test_XMLWriter_1(t *testing.T) {
	var (
		outputfile string          = dataDir("xml_test1.txt")
		result     JUnitTestResult = JUnitTestResult{}
	)

	result.Suites = append(result.Suites, JUnitSuit{Name: "TestPackage", Tests: 10,
		TestCases: []JUnitTestCase{JUnitTestCase{
			ClassName: "TestPackage",
			Name:      "TestCase1",
			Time:      "2.00"}}})

	if err := writeToXML(result, outputfile); err != nil {
		t.Fatal("expected to write XML file, ", err)
	}
	defer os.Remove(outputfile)

	data := []string{
		xml.Header,
		`<testsuites>`,
		`<testsuite name="TestPackage" tests="10" errors="0" failures="0" skip="0">`,
		`<testcase classname="TestPackage" name="TestCase1" time="2.00"></testcase>`,
		`</testsuite>`,
		`</testsuites>`,
	}

	checkXMLOutput(outputfile, data, t)
}

func Test_XMLWriter_2(t *testing.T) {
	var (
		outputfile string          = dataDir("xml_test2.txt")
		result     JUnitTestResult = JUnitTestResult{}
	)

	result.Suites = append(result.Suites,
		JUnitSuit{Name: "TestPackage", Tests: 10,
			TestCases: []JUnitTestCase{JUnitTestCase{
				ClassName: "TestPackage",
				Name:      "TestCase1",
				Time:      "1.00",
				Message: &JUnitTestCaseMessage{
					XMLName: xml.Name{Local: "failure"},
					Message: "message",
					Type:    "type",
					Content: "Content"}}}})

	if err := writeToXML(result, outputfile); err != nil {
		t.Fatal("expected to write XML file, ", err)
	}
	defer os.Remove(outputfile)

	data := []string{
		xml.Header,
		`<testsuites>`,
		`<testsuite name="TestPackage" tests="10" errors="0" failures="0" skip="0">`,
		`<testcase classname="TestPackage" name="TestCase1" time="1.00">`,
		`<failure message="message" type="type">Content</failure>`,
		`</testcase>`,
		`</testsuite>`,
		`</testsuites>`,
	}

	checkXMLOutput(outputfile, data, t)
}

func Test_XMLWriter_3(t *testing.T) {
	var (
		outputfile string          = dataDir("xml_test3.txt")
		result     JUnitTestResult = JUnitTestResult{}
	)

	result.Suites = append(result.Suites,
		JUnitSuit{Name: "TestPackage", Tests: 10, Errors: 3, Failures: 2, Skip: 4,
			TestCases: []JUnitTestCase{JUnitTestCase{
				ClassName: "TestPackage",
				Name:      "TestCase1",
				Time:      "1.00"},
				JUnitTestCase{
					ClassName: "TestPackage",
					Name:      "TestCase2",
					Time:      "1.50",
					Message: &JUnitTestCaseMessage{
						XMLName: xml.Name{Local: "failure"},
						Message: "message",
						Type:    "type",
						Content: "Content"}}}})

	if err := writeToXML(result, outputfile); err != nil {
		t.Fatal("expected to write XML file, ", err)
	}
	defer os.Remove(outputfile)

	data := []string{
		xml.Header,
		`<testsuites>`,
		`<testsuite name="TestPackage" tests="10" errors="3" failures="2" skip="4">`,
		`<testcase classname="TestPackage" name="TestCase1" time="1.00"></testcase>`,
		`<testcase classname="TestPackage" name="TestCase2" time="1.50">`,
		`<failure message="message" type="type">Content</failure>`,
		`</testcase>`,
		`</testsuite>`,
		`</testsuites>`,
	}

	checkXMLOutput(outputfile, data, t)
}

func Test_XMLWriter_4(t *testing.T) {
	var (
		outputfile string          = dataDir("xml_test4.txt")
		result     JUnitTestResult = JUnitTestResult{}
	)

	result.Suites = append(result.Suites, JUnitSuit{Name: "TestPackage1"})
	result.Suites = append(result.Suites,
		JUnitSuit{Name: "TestPackage2", Tests: 10, Errors: 3, Failures: 2, Skip: 4,
			TestCases: []JUnitTestCase{JUnitTestCase{
				ClassName: "TestPackage",
				Name:      "TestCase1",
				Time:      "1.00"},
				JUnitTestCase{
					ClassName: "TestPackage",
					Name:      "TestCase2",
					Time:      "1.50",
					Message: &JUnitTestCaseMessage{
						XMLName: xml.Name{Local: "skipped"},
						Message: "message",
						Type:    "type",
						Content: "Content"}}}})

	if err := writeToXML(result, outputfile); err != nil {
		t.Fatal("expected to write XML file, ", err)
	}
	defer os.Remove(outputfile)

	data := []string{
		xml.Header,
		`<testsuites>`,
		`<testsuite name="TestPackage1" tests="0" errors="0" failures="0" skip="0"></testsuite>`,
		`<testsuite name="TestPackage2" tests="10" errors="3" failures="2" skip="4">`,
		`<testcase classname="TestPackage" name="TestCase1" time="1.00"></testcase>`,
		`<testcase classname="TestPackage" name="TestCase2" time="1.50">`,
		`<skipped message="message" type="type">Content</skipped>`,
		`</testcase>`,
		`</testsuite>`,
		`</testsuites>`,
	}

	checkXMLOutput(outputfile, data, t)
}

func Test_GoTestParser_1(t *testing.T) {
	input := openTestFile(gotestDir("gotest_1.txt"), t)
	defer input.Close()

	result, err := parseGoTest(input)
	if err != nil {
		t.Fatal(err)
	}

	data := `
	<testsuites>
		  <testsuite name="gojunitxml/package_1" tests="1" errors="0" failures="0" skip="0">
		    <testcase classname="gojunitxml/package_1" name="TestPackage_1" time="0.00"></testcase>
		  </testsuite>
	  	  <testsuite name="gojunitxml/package_2" tests="2" errors="0" failures="0" skip="0">
		    <testcase classname="gojunitxml/package_2" name="TestPackage_1" time="0.00"></testcase>
		    <testcase classname="gojunitxml/package_2" name="TestPackage_2" time="0.00"></testcase>
		  </testsuite>
    </testsuites>
    `
	checkTestResult([]byte(data), result, t)
}

func Test_GoTestParser_2(t *testing.T) {
	input := openTestFile(gotestDir("gotest_2.txt"), t)
	defer input.Close()

	result, err := parseGoTest(input)
	if err != nil {
		t.Fatal(err)
	}

	data := `
	<testsuites>
		  <testsuite name="gojunitxml/package_1" tests="2" errors="0" failures="1" skip="1">
		    <testcase classname="gojunitxml/package_1" name="TestPackage_1" time="0.00">
	            <failure type="gotest.error" message="error">package_1_test.go:6: Some error message;</failure>
		    </testcase>
		    <testcase classname="gojunitxml/package_1" name="TestPackage_2" time="0.00">
		  		<skipped type="gotest.skipped" message="skipped">package_1_test.go:10: Some error message,;Some error message in new line;</skipped>
		    </testcase>
		  </testsuite>
	  	  <testsuite name="gojunitxml/package_2" tests="2" errors="0" failures="2" skip="0">
		    <testcase classname="gojunitxml/package_2" name="TestPackage_1" time="0.00">
				<failure type="gotest.error" message="error">package_2_test.go:6: Some error message;</failure>
		    </testcase>
		    <testcase classname="gojunitxml/package_2" name="TestPackage_2" time="0.00">
				<failure type="gotest.error" message="error">package_2_test.go:10: Some error message;</failure>
		    </testcase>
		  </testsuite>
    </testsuites>
    `
	checkTestResult([]byte(data), result, t)
}

func Test_GoTestParser_3(t *testing.T) {
	input := openTestFile(gotestDir("gotest_3.txt"), t)
	defer input.Close()

	result, err := parseGoTest(input)
	if err != nil {
		t.Fatal(err)
	}

	data := `
	<testsuites>
		  <testsuite name="gojunitxml/package_6" tests="0" errors="0" failures="0" skip="0"></testsuite>
		  <testsuite name="gojunitxml/package_1" tests="3" errors="0" failures="2" skip="0">
		    <testcase classname="gojunitxml/package_1" name="TestPackage_1" time="0.00"></testcase>
		    <testcase classname="gojunitxml/package_1" name="TestPackage_2" time="0.00">
	            <failure type="gotest.error" message="error">package_1_test.go:9: Some error message;</failure>
		    </testcase>
		    <testcase classname="gojunitxml/package_1" name="TestPackage_3" time="0.00">
		  		<failure type="gotest.error" message="error">package_1_test.go:13: Some error message,;Some error message in new line;</failure>
		    </testcase>
		  </testsuite>
		  <testsuite name="gojunitxml/package_3" tests="0" errors="0" failures="0" skip="0"></testsuite>
	  	  <testsuite name="gojunitxml/package_2" tests="1" errors="0" failures="0" skip="0">
		    <testcase classname="gojunitxml/package_2" name="TestPackage_1" time="0.00"></testcase>
		  </testsuite>
		  <testsuite name="gojunitxml/package_4" tests="0" errors="0" failures="0" skip="0"></testsuite>
		  <testsuite name="gojunitxml/package_5" tests="0" errors="0" failures="0" skip="0"></testsuite>
    </testsuites>
    `
	checkTestResult([]byte(data), result, t)
}

func Test_GoTestParser_4(t *testing.T) {
	expected := ErrParserError("=== RUN Test$Package_1")
	input := openTestFile(gotestDir("gotest_4.txt"), t)
	defer input.Close()

	if _, err := parseGoTest(input); !reflect.DeepEqual(err, expected) {
		t.Errorf("expected err=%v, got %v", expected, err)
	}
}

func Test_GoTestParser_5(t *testing.T) {
	expected := ErrParserError("ok  	gojunitxml/pac&kage_2	0.006s")
	input := openTestFile(gotestDir("gotest_5.txt"), t)
	defer input.Close()

	if _, err := parseGoTest(input); !reflect.DeepEqual(err, expected) {
		t.Errorf("expected err=%v, got %v", expected, err)
	}
}

func Test_GoTestParser_6(t *testing.T) {
	input := openTestFile(gotestDir("gotest_6.txt"), t)
	defer input.Close()

	result, err := parseGoTest(input)
	if err != nil {
		t.Fatal(err)
	}

	data := `
	<testsuites>
		  <testsuite name="go.junitxml/package_1" tests="1" errors="0" failures="1" skip="0">
		    <testcase classname="go.junitxml/package_1" name="TestPackage_2" time="0.00">
		  		<failure type="gotest.error" message="error">package_1_test.go:10: Some error message,;Some error message in new line;</failure>
		    </testcase>
		  </testsuite>
    </testsuites>
    `
	checkTestResult([]byte(data), result, t)
}

func Test_GoTestParser_7(t *testing.T) {
	input := openTestFile(gotestDir("gotest_7.txt"), t)
	defer input.Close()

	result, err := parseGoTest(input)
	if err != nil {
		t.Fatal(err)
	}

	data := `
	<testsuites>
		  <testsuite name="gojunitxml/package_2" tests="2" errors="0" failures="1" skip="0">
		    <testcase classname="gojunitxml/package_2" name="TestPackage_1" time="0.00">
		  		<failure type="gotest.error" message="error">package_2_test.go:6: Some error message;package_2_test.go:7: Some error message;package_2_test.go:8: Some error message;</failure>
		    </testcase>
		    <testcase classname="gojunitxml/package_2" name="TestPackage_2" time="0.00"></testcase>
		  </testsuite>
    </testsuites>
    `
	checkTestResult([]byte(data), result, t)
}
