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
func byteToBitSlice(byteSlice uint32, length uint8) []bool {
	bits := make([]bool, length)
	var i uint8
	//The length of 32 means we're decoding a 32 bit sequence
	//Loop through the byte and turn it into bit sequence using AND and masking
	//Using an unsigned integer, so
	//7 -> 0
	for i = 0; i < length; i++ {
		mask := byte(1 << i)
		if (byteSlice & uint32(mask)) > 0 {
			bits[(length-1)-i] = true
		} else {
			bits[(length-1)-i] = false
		}

	}
	fmt.Print(bits)
	return bits
}

//Takes a 8 sized bool slice and turns it into a single byte, or 4 bytes if we're converting length
func bitSliceToByte(bitSlice *[]bool, length uint8) []byte {
	var i uint8 = 0
	var j uint8 = 0
	resultingBytes := make([]byte, length)
	for i = 0; i < length; i++ {
		for j = 0; j < 8; j++ {
			if (*bitSlice)[(i*8)+j] {
				resultingBytes[i] |= 1 << (7 - j)
			}
		}
	}

	return resultingBytes
}
