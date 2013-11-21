gojunitxml
==========

`go test -v` output to JUnit XML format.

Installation
------------

    go get github.com/JSzwec/gojunitxml

Usage
-----

    go test -v | gojunitxml -output test_report.xml
    gojunitxml -input gotest_report.txt -output test_report.xml

