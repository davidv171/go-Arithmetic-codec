package main

import (
	"testing"
)

func TestArithmeticDecoder_readFreqTable(t *testing.T) {
	type fields struct {
		inputBits      []bool
		highTable      []uint32
		lowTable       []uint32
		symbolInterval uint32
		step           uint32
		low            uint32
		high           uint32
		output         []byte
	}
	type args struct {
		data []uint8
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Basic short test", fields{}, args{[]uint8{53, 2}}}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arithmeticDecoder := &ArithmeticDecoder{
				inputBits:      tt.fields.inputBits,
				highTable:      tt.fields.highTable,
				lowTable:       tt.fields.lowTable,
				symbolInterval: tt.fields.symbolInterval,
				step:           tt.fields.step,
				low:            tt.fields.low,
				high:           tt.fields.high,
				output:         tt.fields.output,
			}
			arithmeticDecoder.readFreqTable(tt.args.data)
		})
	}
}

func TestArithmeticDecoder_intervalToSymbol(t *testing.T) {
	type fields struct {
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
	type args struct {
		symbolInterval uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   uint8
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arithmeticDecoder := &ArithmeticDecoder{
				inputBits:          tt.fields.inputBits,
				highTable:          tt.fields.highTable,
				lowTable:           tt.fields.lowTable,
				symbolInterval:     tt.fields.symbolInterval,
				step:               tt.fields.step,
				low:                tt.fields.low,
				high:               tt.fields.high,
				output:             tt.fields.output,
				numberOfAllSymbols: tt.fields.numberOfAllSymbols,
				currentInput:       tt.fields.currentInput,
			}
			if got := arithmeticDecoder.intervalToSymbol(tt.args.symbolInterval); got != tt.want {
				t.Errorf("ArithmeticDecoder.intervalToSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}
