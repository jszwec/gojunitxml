gojunitxml
==========

`go test -v` output in JUnit XML format which can be read by eg Jenkins.

Installation
------------

    go get github.com/jszwec/gojunitxml/cmd/gojunitxml

Usage
-----

    go test -v | gojunitxml -output test_report.xml
    gojunitxml -input gotest_report.txt -output test_report.xml
