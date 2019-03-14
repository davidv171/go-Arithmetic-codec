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
		- quarters: 4 sized array containing certain quarters used for error calculation
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
	quarters           []uint32
	index              uint32
}

//The first 256 32-bit bytes are our frequency table, or rather out HIGH table, containing the high value of each symbol
func (arithmeticDecoder *ArithmeticDecoder) readFreqTable(data []uint8) {
	highTable := make([]uint32, 256)
	for i := 4; i < 256*4; i += 4 {
		readBytes := data[i-4 : i]
		convertedInt := bytesToUint32(readBytes)
		highTable[(i-4)/4] = convertedInt

	}
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
	arithmeticDecoder.quarterize(arithmeticDecoder.high)
	arithmeticDecoder.intervalGeneration()
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
	quarters := arithmeticDecoder.quarters
	numberOfAllSymbols := arithmeticDecoder.numberOfAllSymbols
	//An index to keep track of getting the value of the next bit from input bits and adding it to the symbol
	index := arithmeticDecoder.index
	inputBits := arithmeticDecoder.inputBits
	for i := 256 * 4; uint32(i) < 1024+arithmeticDecoder.numberOfAllSymbols; i++ {
		low := arithmeticDecoder.low
		high := arithmeticDecoder.high
		step := arithmeticDecoder.step
		currentBits := arithmeticDecoder.currentInput
		currentByte := arbitraryBitsToByte(&currentBits)
		step = (high - low + 1) / numberOfAllSymbols
		symbolInterval := (uint32(currentByte) - low) / step
		//Calculate the symbol that relates to the symbolInterval
		fmt.Println("CURRENT BYTE", currentByte, " Step ", step, " v ", symbolInterval, " high ", high, " low ", low)
		symbol := arithmeticDecoder.intervalToSymbol(symbolInterval)
		high = low + step*arithmeticDecoder.highTable[symbol] - 1
		low = low + step*arithmeticDecoder.lowTable[symbol]
		fmt.Println("Index ", index, " Field ", currentBits, " Step ", step, " v ", symbolInterval, " symbol: ", symbol, " high ", high, " low ", low)
		//Error intervals
		for (high < quarters[1]) || low >= quarters[1] {
			//E1
			if high < quarters[1] {
				low = 2 * low
				high = 2*high + 1
				//Turn bool into 0 or 1 then add it
				var add uint8 = 0
				if inputBits[index] {
					add = 1
				}
				tempByte := currentByte
				tempByte = 2*tempByte + add
				currentByte = tempByte
				fmt.Println("E1 ", tempByte, " added ", add)
				index++
				//E2

			} else if low >= quarters[1] {
				low = 2 * (low - quarters[1])
				high = 2*(high-quarters[1]) + 1
				var add uint8 = 0
				if inputBits[index] {
					add = 1
				}
				tempByte := currentByte
				tempByte = uint8(2*(uint32(tempByte)-quarters[1]) + uint32(add))
				currentByte = tempByte
				fmt.Println("E2 ", tempByte, " added ", add)
				index++
			}

		}
		for (low >= quarters[0]) && high < quarters[2] {
			low = 2 * (low - quarters[0])
			high = 2*(high-quarters[0]) + 1
			var add uint8 = 0
			if inputBits[index] {
				add = 1
			}
			tempByte := currentByte
			tempByte = uint8(2*(uint32(tempByte)-quarters[0]) + uint32(add))
			currentByte = tempByte
			fmt.Println("E3 ", tempByte, " added ", add)
			index++

		}
		arithmeticDecoder.currentInput = byteToBitSlice(uint32(currentByte), 7)
		arithmeticDecoder.high = high
		arithmeticDecoder.low = low
		arithmeticDecoder.step = step
		arithmeticDecoder.index = index
	}

}

//TODO: Cache the found symbol as it will be the same symbol look up a lot
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

//In case we want to build without a coder
func (arithmeticDecoder *ArithmeticDecoder) quarterize(upperLimit uint32) []uint32 {
	for i := 0; i < 3; i++ {
		// turn to uint64 then back
		upperLimitConverted := uint64(upperLimit)
		arithmeticDecoder.quarters[i] = uint32(((upperLimitConverted + 1) / 4) * uint64(i+1))
	}
	//2^32 or 2^32 - 1 shouldnt matter, we avoid a lot of conversion with it
	arithmeticDecoder.quarters[3] = upperLimit
	return arithmeticDecoder.quarters

}
