package gojunitxml

import (
	"strings"
	"testing"
)

func check(t *testing.T, out, exp string, i int) {
	o, e := strings.Split(out, "\n"), strings.Split(exp, "\n")
	if len(o) != len(e) {
		t.Fatalf("want %s; got %s (i=%d)", exp, out, i)
	}
	for j := range o {
		if a, b := strings.Trim(o[j], "\t\r "), strings.Trim(e[j], "\t\r "); a != b {
			t.Errorf("want %s; got %s (i=%d) (j=%d)", b, a, i, j)
		}
	}
}

func TestParser(t *testing.T) {
	for i := range fixture {
		out, err := Parse(fixture[i].input).Marshal()
		if err != nil {
			t.Fatalf("want err=nil; got %v (i=%d)", err, i)
		}
		check(t, string(out), fixture[i].exp, i)
	}
}
