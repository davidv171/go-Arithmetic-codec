package main

import "fmt"

//Arithmetic coder struct, containing the 1D frequency table , which contains the amount of times each unique symbol repeated
//There will only be one of these used, so we'd probably not need a reflection?
type ArithmeticCoder struct {
	frequencyTable []uint32
	//Table of highs and low, for easier calculations and access, we will be writing the highTable on disk
	highTable []uint32
	lowTable  []uint32
	//Table that keeps track of the order that the bytes have shown up
	readByteSequence []uint8
	//Number of unique symbols in the frequency array, needed for algorithm calculations
	numberOfUniqueSymbols uint8
	//The upper limit of the thing
	upperLimit uint32
}

//Increment the frequency table at a certain index

func (arithmeticCoder *ArithmeticCoder) frequencyTableGenerator(data []uint8) {
	//Keeps track which unique bytes came first when reading from left to right
	//The index at 0 is meant to be the first unique byte we've read
	//Keeps track where in the readByteSequence we've left off, for simpler indexing
	lastIndex := 0
	//Get each byte, increment its frequency
	for i := 0; i < len(data); i++ {
		currentByte := data[i]
		//If the current symbol we're looking at is not already present, means we have a unique symbol
		//TODO: try to optimize this to create a lowTable and highTable out of a single loop
		if arithmeticCoder.frequencyTable[currentByte] == 0 {
			fmt.Print(" ", currentByte, " ")
			arithmeticCoder.numberOfUniqueSymbols++
			arithmeticCoder.readByteSequence[lastIndex] = currentByte
			lastIndex++
		}
		arithmeticCoder.frequencyTable[currentByte]++
	}
	fmt.Print("\n")

	arithmeticCoder.generateHighTable()
}
func (arithmeticCoder *ArithmeticCoder) generateHighTable() {
	var prevHigh uint32 = 0
	var i uint8 = 0
	for i = 0; i < arithmeticCoder.numberOfUniqueSymbols; i++ {
		currSymbol := arithmeticCoder.readByteSequence[i]
		currSymbolFrequency := arithmeticCoder.frequencyTable[currSymbol]
		//We ignore 0 frequencies
		highTableEntry := prevHigh + currSymbolFrequency
		arithmeticCoder.highTable[currSymbol] = highTableEntry
		if currSymbolFrequency != 0 {
			prevHigh = arithmeticCoder.highTable[currSymbol]

		}
	}
}
