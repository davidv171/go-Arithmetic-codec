package main

import (
	"fmt"
	"os"
)

func main() {
	/*1. read the file
	2. calculate the output stream using the arithmetic codec algorithm
	TODO: Frequency table, that is basically a 1D table, index is the value, */
	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	if operation == "c" {
		readBinaryFile(inputFile, operation)

		fmt.Print("Ouputting to... ", outputFile)
	}
}
