package main

import (
	"fmt"
	"os"
)

func main() {
	/*1. read the file
	2. calculate the output stream using the arithmetic codec algorithm
	 Frequency table, that is basically a 1D table, index is the value,

	*/
	operation := os.Args[1]
	inputFile := os.Args[2]
	outputFile := os.Args[3]
	if operation == "c" {
		/*
			TRUE: model creation, because we run the same function twice, once to create a model and once to then run the arithmetic coder")
		*/
		//Empty frequency table, initialized once per all buffer overflows
		frequencyTable := make([]uint32, 256)
		lowTable := make([]uint64, 256)
		highTable := make([]uint64, 256)
		readSequence := make([]uint8, 256)
		quarters := make([]uint64, 4)
		//A series of 0(false) and 1(true) that is then written into bytes and written into the binary compressed file
		//TODO: turn this into output byte array
		outputBits := make([]bool, 0)
		var upperLimit uint64 = 4294967295
		//Initialize an arithmetic codec with empty values except for the upper limit, which has the value of 2^32-1
		//After creating the model is done, we go on to interval creation
		arithmeticCoder := &ArithmeticCoder{frequencyTable, highTable, lowTable, readSequence, 0, upperLimit, 0, upperLimit, 0, 0, quarters, 0, outputBits, 0}
		arithmeticCoder.quarterize(upperLimit)
		readBinaryFile(arithmeticCoder, inputFile, operation, true)
		readBinaryFile(arithmeticCoder, inputFile, operation, false)
		fmt.Print("Ouputting to... ", outputFile)
	}
}
