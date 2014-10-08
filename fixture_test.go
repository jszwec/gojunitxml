package gojunitxml

import (
	"encoding/xml"
	"strings"
)

var fixture = []struct {
	input *strings.Reader
	exp   string
}{
	0: {
		input: strings.NewReader("=== RUN TestPackage_1\n" +
			"--- PASS: TestPackage_1 (0.00 seconds)\n" +
			"PASS\n" +
			"ok    gojunitxml/package_1  0.006s\n" +
			"=== RUN TestPackage_1\n" +
			"--- PASS: TestPackage_1 (0.00 seconds)\n" +
			"=== RUN TestPackage_2\n" +
			"--- PASS: TestPackage_2 (0.00 seconds)\n" +
			"PASS\n" +
			"ok    gojunitxml/package_2  0.005s"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"gojunitxml/package_1.package_1\" tests=\"1\" errors=\"0\" failures=\"0\" skip=\"0\">\n" +
			"<testcase classname=\"gojunitxml/package_1.package_1\" name=\"TestPackage_1\" time=\"0.00\"></testcase>\n" +
			"</testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_2.package_2\" tests=\"2\" errors=\"0\" failures=\"0\" skip=\"0\">\n" +
			"<testcase classname=\"gojunitxml/package_2.package_2\" name=\"TestPackage_1\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"gojunitxml/package_2.package_2\" name=\"TestPackage_2\" time=\"0.00\"></testcase>\n" +
			"</testsuite>\n" +
			"</testsuites>",
	},
	1: {
		input: strings.NewReader("?     gojunitxml/package_6  [no test files]\n" +
			"=== RUN TestPackage_1\n" +
			"--- PASS: TestPackage_1 (0.00 seconds)\n" +
			"=== RUN TestPackage_2\n" +
			"--- SKIP: TestPackage_2 (0.00 seconds)\n" +
			"package_1_test.go:9: Some error message\n" +
			"=== RUN TestPackage_3\n" +
			"--- FAIL: TestPackage_3 (0.00 seconds)\n" +
			"package_1_test.go:13: Some error message,\n" +
			"Some error message in new line\n" +
			"FAIL\n" +
			"FAIL  gojunitxml/package_1  0.005s\n" +
			"?     gojunitxml/package_3  [no test files]\n" +
			"=== RUN TestPackage_1\n" +
			"--- PASS: TestPackage_1 (0.00 seconds)\n" +
			"PASS\n" +
			"ok    gojunitxml/package_2  0.006s\n" +
			"?     gojunitxml/package_4  [no test files]\n" +
			"?     gojunitxml/package_5  [no test files]"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"gojunitxml/package_6.package_6\" tests=\"0\" errors=\"0\" failures=\"0\" skip=\"0\"></testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_1.package_1\" tests=\"3\" errors=\"0\" failures=\"1\" skip=\"1\">\n" +
			"<testcase classname=\"gojunitxml/package_1.package_1\" name=\"TestPackage_1\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"gojunitxml/package_1.package_1\" name=\"TestPackage_2\" time=\"0.00\">\n" +
			"<skipped type=\"gotest.skipped\" message=\"skipped\">package_1_test.go:9: Some error message</skipped>\n" +
			"</testcase>\n" +
			"<testcase classname=\"gojunitxml/package_1.package_1\" name=\"TestPackage_3\" time=\"0.00\">\n" +
			"<failure type=\"gotest.error\" message=\"error\">package_1_test.go:13: Some error message,</failure>\n" +
			"<failure type=\"gotest.error\" message=\"error\">Some error message in new line</failure>\n" +
			"</testcase>\n" +
			"</testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_3.package_3\" tests=\"0\" errors=\"0\" failures=\"0\" skip=\"0\"></testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_2.package_2\" tests=\"1\" errors=\"0\" failures=\"0\" skip=\"0\">\n" +
			"<testcase classname=\"gojunitxml/package_2.package_2\" name=\"TestPackage_1\" time=\"0.00\"></testcase>\n" +
			"</testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_4.package_4\" tests=\"0\" errors=\"0\" failures=\"0\" skip=\"0\"></testsuite>\n" +
			"<testsuite name=\"gojunitxml/package_5.package_5\" tests=\"0\" errors=\"0\" failures=\"0\" skip=\"0\"></testsuite>\n" +
			"</testsuites>",
	},
	2: {
		input: strings.NewReader("=== RUN TestPackage_1\n" +
			"--- FAIL: TestPackage_1 (0.00 seconds)\n" +
			"package_2_test.go:6: Some error message\n" +
			"package_2_test.go:7: Some error message\n" +
			"package_2_test.go:8: Some error message\n" +
			"=== RUN TestPackage_2\n" +
			"--- PASS: TestPackage_2 (0.00 seconds)\n" +
			"FAIL\n" +
			"exit status 1\n" +
			"FAIL  gojunitxml/package_2  0.005s"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"gojunitxml/package_2.package_2\" tests=\"2\" errors=\"0\" failures=\"1\" skip=\"0\">\n" +
			"<testcase classname=\"gojunitxml/package_2.package_2\" name=\"TestPackage_1\" time=\"0.00\">\n" +
			"<failure type=\"gotest.error\" message=\"error\">package_2_test.go:6: Some error message</failure>\n" +
			"<failure type=\"gotest.error\" message=\"error\">package_2_test.go:7: Some error message</failure>\n" +
			"<failure type=\"gotest.error\" message=\"error\">package_2_test.go:8: Some error message</failure>\n" +
			"</testcase>\n" +
			"<testcase classname=\"gojunitxml/package_2.package_2\" name=\"TestPackage_2\" time=\"0.00\"></testcase>\n" +
			"</testsuite>\n" +
			"</testsuites>",
	},
	3: {
		input: strings.NewReader("=== RUN TestHelloWorld\n" +
			"--- PASS: TestHelloWorld (0.00 seconds)\n" +
			"=== RUN TestHelloWorld5\n" +
			"--- PASS: TestHelloWorld5 (0.00 seconds)\n" +
			"=== RUN TestHelloWorld1\n" +
			"=== RUN TestHelloWorld2\n" +
			"printf\n" +
			" " +
			"=== RUN TestHelloWorld3\n" +
			"=== RUN TestHelloWorld4\n" +
			"--- PASS: TestHelloWorld1 (0.00 seconds)\n" +
			"--- FAIL: TestHelloWorld2 (0.00 seconds)\n" +
			"main_test.go:23:\n" +
			"--- PASS: TestHelloWorld3 (0.00 seconds)\n" +
			"printf\n" +
			"--- PASS: TestHelloWorld4 (5.00 seconds)\n" +
			"FAIL\n" +
			"exit status 1\n" +
			"FAIL  check  0.002s"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"check.check\" tests=\"6\" errors=\"0\" failures=\"1\" skip=\"0\">\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld5\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld1\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld2\" time=\"0.00\">\n" +
			"<failure type=\"gotest.error\" message=\"error\">main_test.go:23:</failure>\n" +
			"</testcase>\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld3\" time=\"0.00\"></testcase>\n" +
			"<testcase classname=\"check.check\" name=\"TestHelloWorld4\" time=\"5.00\"></testcase>\n" +
			"</testsuite>\n" +
			"</testsuites>",
	},
	4: {
		input: strings.NewReader("--- FAIL: TestHelloWorld2 (0.00 seconds)\n" +
			"main_test.go:23:\n" +
			"exit status 1\n" +
			"FAIL\n" +
			"FAIL  check/  0.002s"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"check/.check/\" tests=\"1\" errors=\"0\" failures=\"1\" skip=\"0\">\n" +
			"<testcase classname=\"check/.check/\" name=\"TestHelloWorld2\" time=\"0.00\">\n" +
			"<failure type=\"gotest.error\" message=\"error\">main_test.go:23:</failure>\n" +
			"</testcase>\n" +
			"</testsuite>\n" +
			"</testsuites>",
	},
	5: {
		input: strings.NewReader("--- FAIL: TestHelloWorld2 (0.00 seconds)\n" +
			"main_test.go:23:\n" +
			"FAIL"),
		exp: xml.Header + "<testsuites></testsuites>",
	},
	6: {
		input: strings.NewReader("=== RUN TestHelloWorld1\n" +
			"PASS\n" +
			"ok    check  0.005s"),
		exp: xml.Header + "<testsuites>\n" +
			"<testsuite name=\"check.check\" tests=\"0\" errors=\"0\" failures=\"0\" skip=\"0\"></testsuite>\n" +
			"</testsuites>",
	},
}
