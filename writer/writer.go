package main

import (
	utils "ReadOrWriteToFile"
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/*
Girilen string bir listeyi channellar vasıtasıyla dinleyip bir dosyaya yazar.
----------------------------------------------------------------------------------------------------------------
Örnek 1 :
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Input %%%%%%%%%%%%%%%%%%%%%%%%%%%%
2
Hello
World
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Output %%%%%%%%%%%%%%%%%%%%%%%%%%%%
HelloWorld
----------------------------------------------------------------------------------------------------------------
Örnek 2 :
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Input %%%%%%%%%%%%%%%%%%%%%%%%%%%%
5
Lorem
 ipsum
 dolor
 sit
 amet
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Output %%%%%%%%%%%%%%%%%%%%%%%%%%%%
Lorem ipsum dolor sit amet
*/

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)
	stdout, err := os.Create("./result.txt")
	utils.CheckError(err)
	defer stdout.Close()
	writer := bufio.NewWriterSize(stdout, 16*1024*1024)
	fmt.Println("How many input will you enter ?")
	inputArrayCount, err := strconv.ParseInt(strings.TrimSpace(utils.ReadLine(reader)), 10, 64)
	utils.CheckError(err)

	var inputArray []string

	for i := 0; i < int(inputArrayCount); i++ {
		fmt.Printf("Enter your %dth input : \n", i+1)
		inputArrayItem := utils.ReadLine(reader)
		inputArray = append(inputArray, inputArrayItem)
	}

	bytesChannel, doneChannel, errChannel := make(chan []byte), make(chan bool), make(chan error)
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	allocBefore := ms.Alloc

	go writeToFile(bytesChannel, doneChannel, errChannel)

	err = <-errChannel
	if err != nil {
		panic(err)
	}

	for _, b := range inputArray {
		bytesChannel <- []byte(b)
		err := <-errChannel
		if err != nil {
			fmt.Fprintf(writer, "Critical error : %s", err.Error())
			break
		}
	}

	doneChannel <- true
	runtime.ReadMemStats(&ms)
	allocAfter := ms.Alloc
	fmt.Printf("Total memory allocated : %d bytes \n", allocAfter-allocBefore)

	b, err := ioutil.ReadFile(utils.FileName)
	if err == nil {
		fmt.Fprintf(writer, "%s\n", string(b))
	} else {
		fmt.Fprintf(writer, "Critical error : %s", err.Error())
	}

	writer.Flush()
}

func writeToFile(bytesChannel chan []byte, doneChannel chan bool, errChannel chan error) {
	stdout, err := os.Create(utils.FileName)
	if err != nil {
		errChannel <- err
	} else {
		errChannel <- nil
	}
	defer stdout.Close()
	for {
		select {
		case byte := <-bytesChannel:
			n, err := stdout.Write(byte)
			errChannel <- err
			fmt.Printf("wrote %d bytes\n", n)
		case <-doneChannel:
			return
		}
	}
}
