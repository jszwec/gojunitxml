package main

import (
	"flag"
	"fmt"
	"os"
)

var inputfile, outputfile string

func init() {
	const (
		inputopt    = "input"
		inputusage  = "input file"
		outputopt   = "output"
		outputusage = "output file"
	)
	flag.StringVar(&inputfile, inputopt, "", inputusage)
	flag.StringVar(&outputfile, outputopt, "", outputusage)
}

func main() {
	var err error

	flag.Parse()
	input := os.Stdin

	if outputfile == "" {
		flag.PrintDefaults()
		return
	}

	if inputfile != "" {
		if input, err = os.Open(inputfile); err != nil {
			fmt.Println(err)
			return
		}
		defer input.Close()
	}

	results, err := parseGoTest(input)
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = writeToXML(results, outputfile); err != nil {
		fmt.Println(err)
	}
}
