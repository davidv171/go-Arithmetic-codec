package main

import "fmt"

//Arithmetic coder struct, containing the 1D frequency table , which contains the amount of times each unique symbol repeated
//There will only be one of these used, so we'd probably not need a reflection?
type ArithmeticCoder struct{
	frequencyTable []uint32
	//Number of unique symbols in the frequency array, needed for algorithm calculations
	numberOfUniqueSymbols uint32

}
//Increment the frequency table at a certain index
func (arithmeticCoder *ArithmeticCoder) incrementFreqTableElement (index uint8){
	arithmeticCoder.frequencyTable[index] ++
}

func (arithmeticCoder *ArithmeticCoder) uniqueSymbolCount (){
	for i:=0 ; i < 255 ; i++ {
		symbolFrequency := arithmeticCoder.frequencyTable[i]
		if symbolFrequency != 0 {
			fmt.Print(i," ")
			arithmeticCoder.numberOfUniqueSymbols ++
		}
	}
}
func (arithmeticCoder *ArithmeticCoder)frequencyTableGenerator(data []uint8){
	//Get each byte, increment it
	for i:=0 ; i < len(data);i++ {
		currentByte := data[i]
		arithmeticCoder.incrementFreqTableElement(currentByte)
	}

}
