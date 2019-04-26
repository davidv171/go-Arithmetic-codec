package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestByteToBitSlice(t *testing.T) {
	expected := []bool{false, false, false, false, false, false, false, false, false, false, false, false, true, false, false, false}
	actual := byteToBitSlice(8, 16)
	if !reflect.DeepEqual(actual, expected) {
		fmt.Println("Failed test case byteToBitSlice, got: ", actual, " instead of ", expected)
	}
}

func Test_byteToBitSlice(t *testing.T) {
	type args struct {
		byteSlice uint32
		length    uint8
	}
	tests := []struct {
		name string
		args args
		want []bool
	}{
		{"8 bits", args{2, 8},
			[]bool{false, false, false, false, false, false, true, false}},
		{"16 bits", args{2, 16},
			[]bool{false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, true, false}},
		{"32 bits", args{2, 32},
			[]bool{false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, false, false,
				false, false, false, false, false, false, true, false}},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := byteToBitSlice(tt.args.byteSlice, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("byteToBitSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bitSliceToByte(t *testing.T) {
	type args struct {
		bitSlice *[]bool
		length   uint8
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{"1 byte", args{&[]bool{false, false, false, false, false, false, true, false},
			1}, []byte{2}},
		{"2 bytes", args{&[]bool{true, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, true},
			2},
			[]byte{128, 1}},
		{"4 bytes", args{&[]bool{true, false, false, false, false, false, false, false,
			false, false, false, false, false, false, false, true,
			false, false, false, false, false, false, false, true,
			false, false, false, false, false, false, false, true},
			4},
			[]byte{128, 1, 1, 1}},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bitSliceToByte(tt.args.bitSlice, tt.args.length); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("bitSliceToByte() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_bytesToUint32(t *testing.T) {
	type args struct {
		data []uint8
	}
	tests := []struct {
		name string
		args args
		want uint32
	}{
		{"Last byte relevant",
			args{[]uint8{0, 0, 0, 255}},
			255},
		{"First byte relevant and last one",
			args{[]uint8{255, 0, 0, 255}},
			4278190335},
		{"All relevant",
			args{[]uint8{255, 255, 255, 255}},
			4294967295},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := bytesToUint32(tt.args.data); got != tt.want {
				t.Errorf("bytesToUint32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_arbitraryBitsToByte(t *testing.T) {
	type args struct {
		bitSlice *[]bool
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"Basic test", args{&[]bool{true, true, true, true, false, true, false}}, 122},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := arbitraryBitsToByte(tt.args.bitSlice); got != tt.want {
				t.Errorf("sevenBitsToByte() = %v, want %v", got, tt.want)
			}
		})
	}
}
