gojunitxml
==========

`go test -v` output to JUnit XML format.

Installation
------------

    go get github.com/jszwec/gojunitxml

Requirements
============

    Go v1.1 or higher

Usage
-----

    go test -v | gojunitxml -output test_report.xml
    gojunitxml -input gotest_report.txt -output test_report.xml

