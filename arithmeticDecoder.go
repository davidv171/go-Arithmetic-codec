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
		- numberOfAllSymbols: number of all, even non-unique symbols of the encoded file
		- currentInput: the field/array of bools, upon which we od operations on, next bit is taken out of inputBits
	*/
	inputBits          []bool
	highTable          []uint32
	lowTable           []uint32
	symbolInterval     uint32
	step               uint32
	low                uint32
	high               uint32
	output             []byte
	numberOfAllSymbols uint32
	currentInput       []bool
}

//The first 256 32-bit bytes are our frequency table, or rather out HIGH table, containing the high value of each symbol
func (arithmeticDecoder *ArithmeticDecoder) readFreqTable(data []uint8) {
	highTable := make([]uint32, 256)
	for i := 4; i < 256*4; i += 4 {
		readBytes := data[i-4 : i]
		convertedInt := bytesToUint32(readBytes)
		highTable[(i-4)/4] = convertedInt

	}
	fmt.Println("Read high table", highTable)
	arithmeticDecoder.highTable = highTable
	//A temporary variable to avoid ovewriting
	sortedHighTable := make([]uint32, 256)
	copy(sortedHighTable, highTable)
	sort.Slice(sortedHighTable, func(i, j int) bool { return sortedHighTable[i] < sortedHighTable[j] })
	//The largest high(last one in the sorted array) tells us how many non-unique symbols we have
	arithmeticDecoder.numberOfAllSymbols = sortedHighTable[255]
	arithmeticDecoder.step = (arithmeticDecoder.high + 1) / arithmeticDecoder.numberOfAllSymbols
	arithmeticDecoder.generateLowTable(sortedHighTable)
	arithmeticDecoder.initializeField(data)
	arithmeticDecoder.intervalGeneration()
	fmt.Println("High table ", arithmeticDecoder.highTable)
	fmt.Println("Low table ", arithmeticDecoder.lowTable)
	fmt.Println("Read bits ", arithmeticDecoder.inputBits)
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

//Take the first 7 bits(or 32 bits, depending on the arithmetic coder size)
func (arithmeticDecoder *ArithmeticDecoder) initializeField(data []uint8) {
	var i uint32 = 0
	bitSlice := make([]bool, 0)
	fmt.Println(arithmeticDecoder.numberOfAllSymbols)
	//Load every bit, that is not in the model(first 256*4 bytes) onwards
	for i = 256 * 4; i < uint32(len(data)); i++ {
		currByte := data[i]
		bitSlice = append(bitSlice, byteToBitSlice(uint32(currByte), 8)...)
	}
	arithmeticDecoder.inputBits = bitSlice
	//First 7 bits
	arithmeticDecoder.currentInput = bitSlice[0:7]

}

/*
1. read probability table
2. intialize the array of first N bits
3. Decode on every iteration:
	v = (field - low) / step
	step = roundDown(high - low + 1) / n
	high = low + step * high(simbol) - 1
	low = low + step * low(simbol)
4. Check for errors on each iteration:
	while((mHigh < secondQuarter) || (mLow >= secondQuarter))
		- E1: high < secondQuarter
			- low = 2*low
			- high = 2*high + 1
			- field = 2 * field + next bit(0 or 1)
		- E2: low >= 2ndQaurter
			- low = 2*(low - 2ndQuarter)
			- high = 2*(high - 2ndQuarter) + 1
			- field = 2*(field - secondQuarter) + next bit
		- E3: (low>=firstQuarter && high < thirdQuarter)
			- low = 2*(low - firstQuarter)
			- high = 2*(high - firstQuarter)+1
			- field = 2*(field - 1stQuarter) + next bit
*/
func (arithmeticDecoder *ArithmeticDecoder) intervalGeneration() {
	firstBits := arithmeticDecoder.currentInput
	fmt.Println("\nFirst bits\n", firstBits)
	firstByte := arbitraryBitsToByte(&firstBits)
	numberOfAllSymbols := arithmeticDecoder.numberOfAllSymbols
	fmt.Println(firstByte)
	for i := 256 * 4; uint32(i) < 1024+arithmeticDecoder.numberOfAllSymbols; i++ {
		low := arithmeticDecoder.low
		high := arithmeticDecoder.high
		step := arithmeticDecoder.step
		currentBits := arithmeticDecoder.currentInput
		currentByte := arbitraryBitsToByte(&currentBits)
		symbolInterval := (uint32(currentByte) - low) / step
		//Calculate the symbol that relates to the symbolInterval
		symbol := arithmeticDecoder.intervalToSymbol(symbolInterval)
		step = (high - low + 1) / numberOfAllSymbols
		high = low + step*arithmeticDecoder.highTable[symbol] - 1
		fmt.Println("Step", step, " v ", symbolInterval, " symbol: ", symbol, " high ", high, " low ", low)

	}
}

//Returns the symbol that is represented by an interval numbe
func (arithmeticDecoder *ArithmeticDecoder) intervalToSymbol(symbolInterval uint32) uint8 {
	for i := 0; i < 256; i++ {
		//If the interval is found anywhere
		if arithmeticDecoder.highTable[i] > symbolInterval && arithmeticDecoder.lowTable[i] <= symbolInterval {
			return uint8(i)
		}
	}
	return 0
}
