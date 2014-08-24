package main

import (
	"flag"
	"log"
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
		os.Exit(1)
	}

	if inputfile != "" {
		if input, err = os.Open(inputfile); err != nil {
			log.Fatal(err)
		}
		defer input.Close()
	}

	results, err := parseGoTest(input)
	if err != nil {
		log.Fatal(err)
	}
	if err = results.WriteToXML(outputfile); err != nil {
		log.Fatal(err)
	}
}
