package main

import (
	utils "ReadOrWriteToFile"
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

/*
Dosyanın hepsini okumadan, dosyadaki her m karakterin ilk n karakterini okur.
----------------------------------------------------------------------------------------------------------------
Örnek 1 :
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Input %%%%%%%%%%%%%%%%%%%%%%%%%%%%
11
4
Lorem ipsum dolor sit amet
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Output %%%%%%%%%%%%%%%%%%%%%%%%%%%%
Lore
 dol
amet
----------------------------------------------------------------------------------------------------------------
Örnek 2 :
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Input %%%%%%%%%%%%%%%%%%%%%%%%%%%%
1
2
Hello
%%%%%%%%%%%%%%%%%%%%%%%%%%%% Output %%%%%%%%%%%%%%%%%%%%%%%%%%%%
He
el
ll
lo
o
*/

func main() {
	reader := bufio.NewReaderSize(os.Stdin, 16*1024*1024)

	mTemp, err := strconv.ParseInt(strings.TrimSpace(utils.ReadLine(reader)), 10, 64)
	utils.CheckError(err)
	m := int(mTemp)
	nTemp, err := strconv.ParseInt(strings.TrimSpace(utils.ReadLine(reader)), 10, 64)
	utils.CheckError(err)
	n := int(nTemp)

	inputString := utils.ReadLine(reader)
	file, err := os.Create(utils.FileName)
	utils.CheckError(err)
	defer file.Close()
	_, err = file.WriteString(inputString)
	utils.CheckError(err)

	bytesChannel, errChannel := make(chan []byte), make(chan error)

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	allocBefore := ms.Alloc

	go readFile(m, n, bytesChannel, errChannel)
	go func() {
		for {
			select {
			case next := <-bytesChannel:
				fmt.Println(string(next))
			default:
			}
		}
	}()
	err = <-errChannel
	if err == io.EOF {
		runtime.ReadMemStats(&ms)
		allocAfter := ms.Alloc
		fmt.Printf("Total memory allocated : %d bytes \n", allocAfter-allocBefore)
	} else {
		utils.CheckError(err)
	}

}

func readFile(m int, n int, bytesChannel chan []byte, errChannel chan error) {
	f, err := os.Open(utils.FileName)
	if err != nil {
		errChannel <- err
		return
	}
	loopCounter := 1
	for {
		b3 := make([]byte, n)
		_, err = f.Read(b3)
		if err != nil {
			errChannel <- err
			return
		}
		_, err = f.Seek(int64(m*loopCounter), 0)
		if err != nil {
			errChannel <- err
			return
		}
		loopCounter++
		bytesChannel <- b3
	}
}
