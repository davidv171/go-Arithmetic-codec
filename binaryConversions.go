package main

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
func byteToBitSlice(byteSlice *uint8) []bool {
	bits := make([]bool, 8)
	var i uint8
	//Loop through the byte and turn it into bit sequence using AND and masking
	//Using an unsigned integer, so
	//7 -> 0
	for i = 0; i < 8; i++ {
		mask := byte(1 << i)
		if (*byteSlice & mask) > 0 {
			bits[7-i] = true
		} else {
			bits[7-i] = false
		}

	}
	return bits
}
func bitSliceToByte(bitSlice *[]bool) byte {
	var i uint8 = 0
	var resultingByte byte

	for i = 0; i < uint8(len(*bitSlice)); i++ {
		if (*bitSlice)[i] {
			resultingByte |= 1 << (7 - i)
		}
	}
	return resultingByte
}
