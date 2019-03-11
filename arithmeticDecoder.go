package main

import "fmt"

type ArithmeticDecoder struct {
	/*
		- inputBits : the bool array of input bits, that gets shifted for error calculations
		- symbolInterval: is the calculated symbol interval, that falls between a symbol's high and low value in the model definition
		- high table is encoded in the file, low table is derived from the high table
		- step,low,high : step, calculation helper used here for data retention
		- output: the decoded byte array, decompressed file
	*/
	inputBits      []bool
	highTable      []uint32
	lowTable       []uint32
	symbolInterval uint32
	step           uint32
	low            uint32
	high           uint32
	output         []byte
}

//The first 256 32-bit bytes are our frequency table, or rather out HIGH table, containing the high value of each symbol
func (arithmeticDecoder *ArithmeticDecoder) readFreqTable(data []uint8) {
	highTable := make([]uint32, 256)
	for i := 4; i < 256*4; i += 4 {
		if i%4 == 0 {
			readBytes := data[i-4 : i]
			convertedInt := bytesToUint32(readBytes)
			highTable[(i-4)/4] = convertedInt
		}
	}
	fmt.Println(highTable)
	arithmeticDecoder.highTable = highTable
}
