package main

import "fmt"

/*
1. Calculate the frequency table
2. Calculate the symbol high table using that
3. Initialize the algorithm with dividing the interval into 4 quarters
4. Code symbols as they come(left to right when reading)
    - Calculate new borders on each loop, changing up the interval leves
    - Check errors, and based on those errors output certain bits
Because we are reading files with a buffer who isn't large enough usually and gets filled up, we need data persistence between iterations through the buffer
*/
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
	//Number of all symbols, counting repeated symbols, used to calculate the step
	numberOfSymbols uint32
	high            uint32
	low             uint32
	step            uint32
	quarters        []uint32
	e3Counter       uint32
	outputBits      []bool
	//How many bits we're going to write, later used to avoid append functions when writing into file
}

//Increment the frequency table at a certain index

func (arithmeticCoder *ArithmeticCoder) frequencyTableGenerator(data []uint8) {
	//Keeps track which unique bytes came first when reading from left to right
	//The index at 0 is meant to be the first unique byte we've read
	//Keeps track where in the readByteSequence we've left off, for simpler indexing
	lastIndex := 0
	//Get each byte, increment its frequency
	fmt.Println("FREQUENCY TALBE ")
	for i := 0; i < len(data); i++ {
		currentByte := data[i]
		//If the current symbol we're looking at is not already present, means we have a unique symbol
		//Cannot be optimized into a single loop because changing every symbol's high means we'd have to reiterate and recalculate for each new change
		//Makes sure we're not throwing in new data on buffer overflows as well
		if arithmeticCoder.frequencyTable[currentByte] == 0 {
			arithmeticCoder.numberOfUniqueSymbols++
			arithmeticCoder.readByteSequence[lastIndex] = currentByte
			lastIndex++
		}
		arithmeticCoder.numberOfSymbols++
		arithmeticCoder.frequencyTable[currentByte]++
	}
	fmt.Println(arithmeticCoder.frequencyTable)
	arithmeticCoder.generateHighTable()
}

//Generate a table of highs for each symbol
//Needs a helper array to keep track of which symbols appeared first in a left-to-right order
func (arithmeticCoder *ArithmeticCoder) generateHighTable() {
	prevHigh := uint32(0)
	//Highest amount the iterator can reach is 255, which is the max amount of unique symbols
	//To save a few iterations, we keep the number of unique symbols
	var i uint8 = 0
	for i = 0; i < arithmeticCoder.numberOfUniqueSymbols; i++ {
		currSymbol := arithmeticCoder.readByteSequence[i]
		currSymbolFrequency := arithmeticCoder.frequencyTable[currSymbol]
		//We ignore 0 frequencies
		highTableEntry := prevHigh + currSymbolFrequency
		lowTableEntry := prevHigh
		arithmeticCoder.highTable[currSymbol] = highTableEntry
		arithmeticCoder.lowTable[currSymbol] = lowTableEntry
		if currSymbolFrequency != 0 {
			prevHigh = highTableEntry

		}
	}
	fmt.Println(arithmeticCoder.highTable)
	fmt.Println(arithmeticCoder.lowTable)
}

//Takes a number and sets up quarters from it
//An array of quarters is then used for algorithm calculation of border changes
func (arithmeticCoder *ArithmeticCoder) quarterize(upperLimit uint32) []uint32 {
	for i := 0; i < 3; i++ {
		// turn to uint64 then back
		upperLimitConverted := uint64(upperLimit)
		arithmeticCoder.quarters[i] = uint32(((upperLimitConverted + 1) / 4) * uint64(i+1))
	}
	//2^32 or 2^32 - 1 shouldnt matter, we avoid a lot of conversion with it
	arithmeticCoder.quarters[3] = upperLimit
	return arithmeticCoder.quarters

}

/*
We use the same file
Calculate the intervals on each run based on the read bytes:
    - step = roundDown(high - low + 1)/n
    - high = low + step * high(symbol) - 1
    - low = low + step * low(symbol)
    n = number of all unique symbols
Based on those, we calculate the errors(E1/E2 and or E3)
- E1: (high < 2ndQuarter) -> low = low * 2, high = high * 2 + 1	OUTPUT: 0 and E3_COUNTER times 1, set E3_COUNTER = 0
- E2: (low >= 2ndQuarter) -> low = 2(*low - 2ndQuarter) , high = 2*(high - 2ndQuarter) + 1, OUTPUT: 1 and E3_COUNTER times 0, set E3_COUNTER = 0
- E3: (low >= firstQuarter && high < 3rdQuarter) -> low = 2*(low - 1stQuarter) , high = 2*(high - 1stQuarter) + 1 ... E3_COUNTER++
Every error check gets repeated:
- E1 and E2: while(( high < 2ndQuarter) || low >= 2ndQuarter)
- E3: while(firstQuarter <= low) && (high < i3rdQuarter)
If E1 was calculated, then E2 cannot be, E3 can be calculated no matter the previous error detected
*/

func (arithmeticCoder *ArithmeticCoder) intervalCalculation(data []uint8) {
	high := arithmeticCoder.high
	low := arithmeticCoder.low
	quarters := arithmeticCoder.quarters
	step := arithmeticCoder.step
	e3Counter := arithmeticCoder.e3Counter
	outputBits := make([]bool, 0)
	var i uint64 = 0
	for ; i < uint64(len(data)); i++ {
		step = uint32((uint64(high) - uint64(low) + 1) / uint64(arithmeticCoder.numberOfSymbols))
		high = low + step*arithmeticCoder.highTable[data[i]] - 1
		low = low + step*arithmeticCoder.lowTable[data[i]]
		for (high < quarters[1]) || (low >= quarters[1]) {
			if high < quarters[1] {
				low = low * 2
				high = high*2 + 1
				var j uint32
				//OUTPUT: 0
				outputBits = append(outputBits, false)
				for j = 0; j < e3Counter; j++ {
					//OUTPUT: 1
					outputBits = append(outputBits, true)
				}
				e3Counter = 0
			} else if low >= quarters[1] {
				low = 2 * (low - quarters[1])
				high = 2*(high-quarters[1]) + 1
				//fmt.Print("1")
				outputBits = append(outputBits, true)
				var j uint32
				for j = 0; j < e3Counter; j++ {
					//fmt.Print("0 ")
					outputBits = append(outputBits, false)

				}
				e3Counter = 0
			}
		}
		for (quarters[0] <= low) && (high < quarters[2]) {
			if low >= quarters[0] {
				low = 2 * (low - quarters[0])
				high = 2*(high-quarters[0]) + 1
				e3Counter++
			}
		}

	}
	arithmeticCoder.high = high
	arithmeticCoder.low = low
	arithmeticCoder.step = step
	arithmeticCoder.e3Counter = e3Counter
	arithmeticCoder.outputBits = outputBits
}
