package main

type ArithmeticDecoder struct {
	/*
		- inputBits : the bool array of input bits, that gets shifted for error calculations
		- symbolInterval: is the calculated symbol interval, that falls between a symbol's high and low value in the model definition
		- step,low,high : step, calculation helper used here for data retention
		- output: the decoded byte array, decompressed file
	*/
	inputBits      []bool
	symbolInterval uint32
	step           uint32
	low            uint32
	high           uint32
	output         []byte
}

//The first 256 32-bit bytes are our frequency table, or rather out HIGH table, containing the high value of each symbol
func readFreqTable(data []byte) {

}
