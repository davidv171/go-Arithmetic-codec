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
	fmt.Print(fileSize)
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
	frequencyTable := make([]uint32,256)
	arithmeticCoder := ArithmeticCoder{frequencyTable,0}
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
	arithmeticCoder.uniqueSymbolCount()
	fmt.Println("\nUNIQUE : ",arithmeticCoder.numberOfUniqueSymbols)

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