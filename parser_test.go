package gojunitxml

import (
	"encoding/xml"
	"strings"
	"testing"
)

func check(t *testing.T, out, exp string, i int) {
	o, e := strings.Split(out, "\n"), strings.Split(exp, "\n")
	if len(o) != len(e) {
		t.Fatalf("want %s; got %s (i=%d)", exp, out, i)
	}
	for j := range o {
		if a, b := strings.Trim(o[j], "\t\b\r "), strings.Trim(e[j], "\t\b\r "); a != b {
			t.Errorf("want %s; got %s (i=%d) (j=%d)", b, a, i, j)
		}
	}
}

func TestParser(t *testing.T) {
	data := []struct {
		input *strings.Reader
		xml   string
	}{
		0: {
			input: strings.NewReader(`=== RUN TestPackage_1
			--- PASS: TestPackage_1 (0.00 seconds)
			PASS
			ok  	gojunitxml/package_1	0.006s
			=== RUN TestPackage_1
			--- PASS: TestPackage_1 (0.00 seconds)
			=== RUN TestPackage_2
			--- PASS: TestPackage_2 (0.00 seconds)
			PASS
			ok  	gojunitxml/package_2	0.005s`),
			xml: xml.Header + `<testsuites>
			  <testsuite name="gojunitxml/package_1.package_1" tests="1" errors="0" failures="0" skip="0">
			    <testcase classname="gojunitxml/package_1.package_1" name="TestPackage_1" time="0.00"></testcase>
			  </testsuite>
		  	  <testsuite name="gojunitxml/package_2.package_2" tests="2" errors="0" failures="0" skip="0">
			    <testcase classname="gojunitxml/package_2.package_2" name="TestPackage_1" time="0.00"></testcase>
			    <testcase classname="gojunitxml/package_2.package_2" name="TestPackage_2" time="0.00"></testcase>
			  </testsuite>
	    </testsuites>`,
		},
		1: {
			input: strings.NewReader(`?   	gojunitxml/package_6	[no test files]
			=== RUN TestPackage_1
			--- PASS: TestPackage_1 (0.00 seconds)
			=== RUN TestPackage_2
			--- SKIP: TestPackage_2 (0.00 seconds)
				package_1_test.go:9: Some error message
			=== RUN TestPackage_3
			--- FAIL: TestPackage_3 (0.00 seconds)
				package_1_test.go:13: Some error message,
					Some error message in new line
			FAIL
			FAIL	gojunitxml/package_1	0.005s
			?   	gojunitxml/package_3	[no test files]
			=== RUN TestPackage_1
			--- PASS: TestPackage_1 (0.00 seconds)
			PASS
			ok  	gojunitxml/package_2	0.006s
			?   	gojunitxml/package_4	[no test files]
			?   	gojunitxml/package_5	[no test files]`),
			xml: xml.Header + `<testsuites>
				<testsuite name="gojunitxml/package_6.package_6" tests="0" errors="0" failures="0" skip="0"></testsuite>
				<testsuite name="gojunitxml/package_1.package_1" tests="3" errors="0" failures="1" skip="1">
					<testcase classname="gojunitxml/package_1.package_1" name="TestPackage_1" time="0.00"></testcase>
					<testcase classname="gojunitxml/package_1.package_1" name="TestPackage_2" time="0.00">
								<skipped type="gotest.skipped" message="skipped">package_1_test.go:9: Some error message</skipped>
					</testcase>
					<testcase classname="gojunitxml/package_1.package_1" name="TestPackage_3" time="0.00">
						<failure type="gotest.error" message="error">package_1_test.go:13: Some error message,</failure>
						<failure type="gotest.error" message="error">Some error message in new line</failure>
					</testcase>
				</testsuite>
				<testsuite name="gojunitxml/package_3.package_3" tests="0" errors="0" failures="0" skip="0"></testsuite>
					<testsuite name="gojunitxml/package_2.package_2" tests="1" errors="0" failures="0" skip="0">
					<testcase classname="gojunitxml/package_2.package_2" name="TestPackage_1" time="0.00"></testcase>
				</testsuite>
				<testsuite name="gojunitxml/package_4.package_4" tests="0" errors="0" failures="0" skip="0"></testsuite>
				<testsuite name="gojunitxml/package_5.package_5" tests="0" errors="0" failures="0" skip="0"></testsuite>
			</testsuites>`,
		},
		2: {
			input: strings.NewReader(`=== RUN TestPackage_1
			--- FAIL: TestPackage_1 (0.00 seconds)
					package_2_test.go:6: Some error message
					package_2_test.go:7: Some error message
					package_2_test.go:8: Some error message
			=== RUN TestPackage_2
			--- PASS: TestPackage_2 (0.00 seconds)
			FAIL
			exit status 1
			FAIL	gojunitxml/package_2	0.005s`),
			xml: xml.Header + `<testsuites>
				<testsuite name="gojunitxml/package_2.package_2" tests="2" errors="0" failures="1" skip="0">
					<testcase classname="gojunitxml/package_2.package_2" name="TestPackage_1" time="0.00">
		        <failure type="gotest.error" message="error">package_2_test.go:6: Some error message</failure>
						<failure type="gotest.error" message="error">package_2_test.go:7: Some error message</failure>
						<failure type="gotest.error" message="error">package_2_test.go:8: Some error message</failure>
					</testcase>
					<testcase classname="gojunitxml/package_2.package_2" name="TestPackage_2" time="0.00"></testcase>
				</testsuite>
			</testsuites>`,
		},
		3: {
			input: strings.NewReader(`=== RUN TestHelloWorld
			--- PASS: TestHelloWorld (0.00 seconds)
			=== RUN TestHelloWorld5
			--- PASS: TestHelloWorld5 (0.00 seconds)
			=== RUN TestHelloWorld1
			=== RUN TestHelloWorld2
			printf

			=== RUN TestHelloWorld3
			=== RUN TestHelloWorld4
			--- PASS: TestHelloWorld1 (0.00 seconds)
			--- FAIL: TestHelloWorld2 (0.00 seconds)
				main_test.go:23:
			--- PASS: TestHelloWorld3 (0.00 seconds)
			printf
			--- PASS: TestHelloWorld4 (5.00 seconds)
			FAIL
			exit status 1
			FAIL	check	0.002s`),
			xml: xml.Header + `<testsuites>
				<testsuite name="check.check" tests="6" errors="0" failures="1" skip="0">
					<testcase classname="check.check" name="TestHelloWorld" time="0.00"></testcase>
					<testcase classname="check.check" name="TestHelloWorld5" time="0.00"></testcase>
					<testcase classname="check.check" name="TestHelloWorld1" time="0.00"></testcase>
					<testcase classname="check.check" name="TestHelloWorld2" time="0.00">
						<failure type="gotest.error" message="error">main_test.go:23:</failure>
					</testcase>
					<testcase classname="check.check" name="TestHelloWorld3" time="0.00"></testcase>
					<testcase classname="check.check" name="TestHelloWorld4" time="5.00"></testcase>
				</testsuite>
			</testsuites>`,
		},
		4: {
			input: strings.NewReader(`--- FAIL: TestHelloWorld2 (0.00 seconds)
				main_test.go:23:
				exit status 1
			FAIL
			FAIL	check/	0.002s`),
			xml: xml.Header + `<testsuites>
				<testsuite name="check/.check/" tests="1" errors="0" failures="1" skip="0">
				<testcase classname="check/.check/" name="TestHelloWorld2" time="0.00">
					<failure type="gotest.error" message="error">main_test.go:23:</failure>
				</testcase>
				</testsuite>
			</testsuites>`,
		},
		5: {
			input: strings.NewReader(`--- FAIL: TestHelloWorld2 (0.00 seconds)
				main_test.go:23:
			FAIL`),
			xml: xml.Header + `<testsuites></testsuites>`,
		},
		6: {
			input: strings.NewReader(`=== RUN TestHelloWorld1
				PASS
				ok  	check	0.005s`),
			xml: xml.Header + `<testsuites>
			<testsuite name="check.check" tests="0" errors="0" failures="0" skip="0"></testsuite>
			</testsuites>`,
		},
	}
	for i := range data {
		out, err := Parse(data[i].input).Marshal()
		if err != nil {
			t.Fatalf("want err=nil; got %v (i=%d)", err, i)
		}
		check(t, string(out), data[i].xml, i)
	}
}
