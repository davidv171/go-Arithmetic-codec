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
	//4294967295
	var upperLimit uint32 = 4294967295
	quarters := make([]uint32, 4)
	//Compression
	if operation == "c" {
		/*
			TRUE: model creation, because we run the same function twice, once to create a model and once to then run the arithmetic coder")
		*/
		//Empty frequency table, initialized once per all buffer overflows
		frequencyTable := make([]uint32, 256)
		lowTable := make([]uint32, 256)
		highTable := make([]uint32, 256)
		readSequence := make([]uint8, 256)
		//A series of 0(false) and 1(true) that is then written into bytes and written into the binary compressed file
		//TODO: turn this into output byte array
		outputBits := make([]bool, 0)
		//4294967295
		//Initialize an arithmetic codec with empty values except for the upper limit, which has the value of 2^32-1
		//After creating the model is done, we go on to interval creation
		arithmeticCoder := &ArithmeticCoder{frequencyTable, highTable, lowTable,
			readSequence, 0, upperLimit, 0, upperLimit,
			0, 0, quarters, 0, outputBits}
		quarters = arithmeticCoder.quarterize(upperLimit)
		//The last argument is for the arithmetic decoder, whenever we are not decoding, it's nil
		readBinaryFile(arithmeticCoder, inputFile, operation, true, nil, outputFile)
		readBinaryFile(arithmeticCoder, inputFile, operation, false, nil, outputFile)
		fmt.Print("Ouputting to... ", outputFile)
		//Decomepression, read compressed file, deconstruct the symbols based off of it
	} else if operation == "d" {
		/*	inputBits      []bool
			highTable []uint32
			lowTable []uint32
			symbolInterval uint32
			step           uint32
			low            uint32
			high           uint32
			output         []byte
		}*/
		inputBits := make([]bool, 0)
		highTable := make([]uint32, 0)
		lowTable := make([]uint32, 0)
		var symbolInterval uint32
		output := make([]byte, 0)
		currentInputBits := make([]bool, 0)
		arithmeticDecoder := &ArithmeticDecoder{inputBits, highTable, lowTable,
			symbolInterval, 0, 0, upperLimit, output, 0,
			currentInputBits, quarters, 32,nil,nil,nil,inputFile}
		readBinaryFile(nil, inputFile, operation, false, arithmeticDecoder, outputFile)
	}
}
