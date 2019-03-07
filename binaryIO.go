package main

import (
	"fmt"
	"io"
	"os"
)

func readBinaryFile(filepath string, operation string) {
	file, err := os.Open(filepath)
	defer file.Close()
	fileInfo, err := file.Stat()
	errCheck(err)
	fileSize := fileInfo.Size()
	fmt.Print("File size is ")
	fmt.Print(fileSize, "\n")
	//Make sure the search slice fits, this will be our b
	//
	//uffer
	var bufferSize int64
	bufferSize = 4096
	if fileSize < bufferSize {
		bufferSize = fileSize
	}
	var bufferOverflow int64 = 0
	//Data where we put the read bytes into
	data := make([]byte, bufferSize)
	//Empty frequency table, initialized once per all buffer overflows
	frequencyTable := make([]uint32, 256)
	lowTable := make([]uint32, 256)
	highTable := make([]uint32, 256)
	readSequence := make([]uint8, 256)
	var upperLimit uint32 = 4294967295
	//Initialize an arithmetic codec with empty values except for the upper limit, which has the value of 2^32-1
	arithmeticCoder := &ArithmeticCoder{frequencyTable, highTable, lowTable, readSequence, 0, upperLimit}
	for {
		//Loop through the file, retrieve the bytes as integers
		//:cap(data) = capacity of array, how many elements it can take before it has to resize
		//Init slice
		data = data[:cap(data)]
		//Byte in the file
		readByte, err := file.Read(data)
		if err != nil {
			if err == io.EOF {
				fmt.Print("\n")
				fmt.Println("Done reading file")
				break
			}
			fmt.Println(err)
			return
		}
		data = data[:readByte]
		//for _,aByte := range data {
		//Add the new 8 booleans to the end of the bits array
		//bits = append(bits,byteToBitSlice(&aByte)...)
		//}
		if operation == "c" {
			arithmeticCoder.frequencyTableGenerator(data)
		}
		bufferOverflow += bufferSize
		//So we're aware of indexes if the file is larger
	}
	fmt.Println("\nUNIQUE : ", arithmeticCoder.numberOfUniqueSymbols)
	fmt.Print("\nFrequences : ", arithmeticCoder.frequencyTable)
	fmt.Print("\nHigh table : ", arithmeticCoder.highTable, "\n", len(arithmeticCoder.highTable), "\n")

}
func writeBinaryFile(fileName string, bytesToWrite *[]byte, bufferOverflow int64) {
	if bufferOverflow == 0 {
		_, err := os.Create(fileName)
		errCheck(err)
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0644)
	defer file.Close()
	errCheck(err)
	_, err = file.WriteAt(*bytesToWrite, bufferOverflow)
	errCheck(err)
	os.Exit(0)
}
