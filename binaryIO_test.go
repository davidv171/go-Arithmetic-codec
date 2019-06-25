package main

import "testing"

func Test_readBinaryFile(t *testing.T) {
	type args struct {
		arithmeticCoder   *ArithmeticCoder
		filepath          string
		operation         string
		modelCreation     bool
		arithmeticDecoder *ArithmeticDecoder
		outputFile        string
	}
	//4294967295
	inputFile := "output.COMPRESSEDTIME"
	var upperLimit uint32 = 4294967295
	quarters := make([]uint32, 4)
	inputBits := make([]bool, 0)
	highTable := make([]uint32, 0)
	lowTable := make([]uint32, 0)
	var symbolInterval uint32
	output := make([]byte, 0)
	currentInputBits := make([]bool, 0)
	arithmeticDecoder := &ArithmeticDecoder{inputBits, highTable, lowTable,
		symbolInterval, 0, 0, upperLimit, output, 0,
		currentInputBits, quarters, 32, nil, nil, nil, nil, inputFile}
	tests := []struct {
		name string
		args args
	}{

		{"Basic test with lena decompress", args{nil,
			inputFile,"d",false,
			arithmeticDecoder,"output.original"}},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readBinaryFile(tt.args.arithmeticCoder, tt.args.filepath, tt.args.operation, tt.args.modelCreation, tt.args.arithmeticDecoder, tt.args.outputFile)
		})
	}
}
func BenchmarkLena(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		//4294967295
		inputFile := "output.COMPRESSEDTIME"
		var upperLimit uint32 = 4294967295
		quarters := make([]uint32, 4)
		inputBits := make([]bool, 0)
		highTable := make([]uint32, 0)
		lowTable := make([]uint32, 0)
		var symbolInterval uint32
		output := make([]byte, 0)
		currentInputBits := make([]bool, 0)
		arithmeticDecoder := &ArithmeticDecoder{inputBits, highTable, lowTable,
			symbolInterval, 0, 0, upperLimit, output, 0,
			currentInputBits, quarters, 32, nil, nil, nil, nil, inputFile}
		readBinaryFile(nil, inputFile, "d", false, arithmeticDecoder, "backToOriginal")
	}
}
