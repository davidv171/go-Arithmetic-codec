package main

import "fmt"

//Check a string and return its bool slice
//FALSE = binary 0, TRUE = binary 1
func argumentToBinary(argument string) []bool {
	stringLength := len(argument)
	bits := make([]bool, stringLength)
	for position, char := range argument {
		if char == '0' {
			bits[position] = false
		} else if char == '1' {
			bits[position] = true
		}
	}
	return bits
}
func byteToBitSlice(bytes uint32, length uint8) []bool {
	bits := make([]bool, length)
	var i uint8
	//The length of 32 means we're decoding a 32 bit sequence
	//Loop through the and turn it into bit sequence using AND and masking
	//Using an unsigned integer, so
	//7 -> 0
	for i = 0; i < length; i++ {
		mask := uint8(1 << i)
		if (bytes & uint32(mask)) > 0 {
			bits[(length-1)-i] = true
		} else {
			bits[(length-1)-i] = false
		}

	}
	fmt.Print(bits)
	return bits
}

//Takes a length*8 sized bool slice and turns it into a single byte, or 4 bytes if we're converting 32 byte length bitslice
func bitSliceToByte(bitSlice *[]bool, length uint8) []uint8 {
	var i uint8 = 0
	var j uint8 = 0
	resultingBytes := make([]uint8, length)
	for i = 0; i < length; i++ {
		for j = 0; j < 8; j++ {
			if (*bitSlice)[(i*8)+j] {
				resultingBytes[i] |= 1 << (7 - j)
			}
		}
	}

	return resultingBytes
}

//Only for testing purposes, takes 7 bits at a time to generate a byte out of it
func arbitraryBitsToByte(bitSlice *[]bool) uint8 {
	var resultingByte uint8 = 0
	var i uint8 = 0
	length := len((*bitSlice))
	for i = 0; i < uint8(length); i++ {
		if (*bitSlice)[i] {
			resultingByte |= 1 << ((uint8(length - 1)) - i)
		}
	}
	return resultingByte
}

//Turn 4 bytes(uint8) into an uint32 by doing bit shifting magic
//Not done in a loop because it kept messing up for no reason.
func bytesToUint32(data []uint8) uint32 {
	var returnValue uint32 = 0
	returnValue |= uint32(data[3])
	returnValue |= uint32(data[2]) << 8
	returnValue |= uint32(data[1]) << 16
	returnValue |= uint32(data[0]) << 24
	return returnValue
}
