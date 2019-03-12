package main

import (
	"fmt"
	"sort"
)

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
	fmt.Println("Read high table", highTable)
	arithmeticDecoder.highTable = highTable
	sortedHighTable := make([]uint32, 256)
	copy(sortedHighTable, highTable)
	sort.Slice(sortedHighTable, func(i, j int) bool { return sortedHighTable[i] < sortedHighTable[j] })
	arithmeticDecoder.generateLowTable(sortedHighTable)
	fmt.Println("High table", arithmeticDecoder.highTable)
	fmt.Println("Low table", arithmeticDecoder.lowTable)
	//Sort the array so we can generate low table much easier
}

/*Generate low table out of the high table
1. Find the lowest high
2. The lowest high's symbol's low = 0
3. Find the second lowest high, that symbol's low is the previous symbol's high
* This could also be used for interval definition
*/
func (arithmeticDecoder *ArithmeticDecoder) generateLowTable(sortedHighTable []uint32) {
	//Generate low table based on the sorted table and high table that is unsorted
	//Put the value of non-0 indexes of sorted table into the index corresponding to the value in hightable
	//IN other words, put the low table's values into the right
	highTable := arithmeticDecoder.highTable
	lowTable := make([]uint32, 256)
	//Sortedhightable is sorted from lowest to highest
	//The first low value of a symbol is always 0
	//The next value of low is equal to the next high value
	for i := 0; i < 256; i++ {
		currSorted := sortedHighTable[i]
		if currSorted != 0 {
			//Find the index with value in unsorted table, don't loop from start every time
			for j := 0; j < 256; j++ {
				if currSorted == highTable[j] {
					lowTable[j] = sortedHighTable[i-1]
					break
				}
			}
		}
	}
	arithmeticDecoder.lowTable = lowTable

}
