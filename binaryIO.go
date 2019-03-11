package main

import (
	"fmt"
	"io"
	"os"
)

func readBinaryFile(arithmeticCoder *ArithmeticCoder, filepath string, operation string, modelCreation bool) {
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
		}
		bufferOverflow += bufferSize
		//So we're aware of indexes if the file is larger
	}
	/*After we are done encoding the values byte by byte, we look at the rest:
	- if low < firstQuarter : output 01 and E3_COUNTER times bit 1
	- else : output 10 and E3_COUNTER times bit 0
	*/
	if !modelCreation {
		//fmt.Println("The rest:")
		if arithmeticCoder.low < arithmeticCoder.quarters[0] {
			fmt.Print("\n01 ")
			arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, false, true)
			arithmeticCoder.writtenSize++
			for i := 0; uint32(i) < arithmeticCoder.e3Counter; i++ {
				fmt.Print("1")
				arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, true)
				arithmeticCoder.writtenSize++
			}
		} else {
			fmt.Println("10")
			arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, true, false)

			for i := 0; uint32(i) < arithmeticCoder.e3Counter; i++ {
				arithmeticCoder.writtenSize++
				fmt.Print("0")
				arithmeticCoder.outputBits = append(arithmeticCoder.outputBits, false)
			}
		}
		fmt.Println("")
		//Turn the output bits into bytes
		fmt.Println("END OF THING SIZE ", len(arithmeticCoder.outputBits), arithmeticCoder.outputBits)
		var outputBytes []byte
		for i := 0; i < len(arithmeticCoder.outputBits); i += 8 {
			tempSlice := arithmeticCoder.outputBits[i : i+8]
			outputBytes = append(outputBytes, bitSliceToByte(&tempSlice))
		}
		writeBinaryFile("out", &outputBytes, 0)
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
