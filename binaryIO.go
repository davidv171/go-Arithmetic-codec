package main

import (
	"fmt"
	"io"
	"os"
)

func readBinaryFile(arithmeticCoder *ArithmeticCoder, filepath string, operation string, modelCreation bool, arithmeticDecoder *ArithmeticDecoder, outputFile string) {
	file, err := os.Open(filepath)
	defer file.Close()
	fileInfo, err := file.Stat()
	errCheck(err)
	fileSize := fileInfo.Size()
	fmt.Print("File size is ")
	fmt.Print(fileSize, "\n")
	var bufferSize int64
	bufferSize = 4096
	//YOLO
	if fileSize < bufferSize || operation == "d" {
		bufferSize = fileSize
	}
	var bufferOverflow int64 = 0
	//Data where we put the read bytes into
	data := make([]byte, bufferSize)
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
			if modelCreation {
				arithmeticCoder.frequencyTableGenerator(data)
			} else {
				arithmeticCoder.intervalCalculation(data)

			}
		} else if operation == "d" {
			arithmeticDecoder.readFreqTable(data)
			outputBytes := arithmeticDecoder.output
			writeBinaryFile(outputFile, &outputBytes, 0)

		}
		bufferOverflow += bufferSize
		//So we're aware of indexes if the file is larger
	}
	/*After we are done encoding the values byte by byte, we look at the rest:
	- if low < firstQuarter : output 01 and E3_COUNTER times bit 1
	- else : output 10 and E3_COUNTER times bit 0
	*/
	if !modelCreation && arithmeticCoder != nil {
		writeEncoded(arithmeticCoder, outputFile)
		//fmt.Println("The rest:")

	}

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
func writeEncoded(arithmeticCoder *ArithmeticCoder, fileName string) {
	if arithmeticCoder.low < arithmeticCoder.quarters[0] {
		arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, false, true)
		for i := 0; uint32(i) < arithmeticCoder.e3Counter; i++ {
			arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, true)
		}
	} else {
		arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, true, false)

		for i := 0; uint32(i) < arithmeticCoder.e3Counter; i++ {
			arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, false)
		}
	}
	fmt.Println("")
	//Write the 32uint[256] high table into file
	//If the value in high table is 0,
	// we can just write 4 0 bytes into the table, this saves us some time when doing compression on a small amount of unique symbols
	outputBytes := make([]byte, 0)
	//TODO: decide on number of written bytes based on the highest value
	for i := 0; i < 256; i++ {
		currentElement := arithmeticCoder.highTable[i]
		if currentElement != 0 {
			tempSlice := byteToBitSlice(currentElement, 32)
			outputBytes = append(outputBytes, bitSliceToByte(&tempSlice, 4)...)

		} else {
			outputBytes = append(outputBytes, 0, 0, 0, 0)
		}

	}
	for i := 0; i < len(arithmeticCoder.outputBits); i += 8 {
		tempSlice := arithmeticCoder.outputBits[i : i+8]
		outputBytes = append(outputBytes, bitSliceToByte(&tempSlice, 1)[0])
	}
	fmt.Println("Output size ", len(outputBytes))
	writeBinaryFile(fileName, &outputBytes, 0)
}
