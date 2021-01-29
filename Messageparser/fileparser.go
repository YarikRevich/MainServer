package Messageparser

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strings"
	"strconv"
	"os"

	ascii "github.com/galsondor/go-ascii"
)


type FileParser interface{
	GetFile()string
	ParseFile(string)[]*net.UDPAddr
}

type FP struct{}


func(f *FP) GetFile()string{
	file, err := os.OpenFile("server.txt", os.O_CREATE|os.O_RDONLY, 0666)
	if err != nil{
		panic(err)
	}
	defer file.Close()
	var buff bytes.Buffer
	io.Copy(&buff, file)

	if buff.Len() == 0{
		fmt.Println("The configuration file is empty!")
		os.Exit(0)
	}

	var cleaned string
	for _, symbol := range buff.String(){
		if ascii.IsGraph(byte(symbol)) || ascii.IsSpace(byte(symbol)){
			cleaned += string(symbol)
		}
	}
	return cleaned
}

func(f *FP) ParseFile(info string)[]*net.UDPAddr{
	replacedInfo := strings.Replace(info, "\n", " ", -1)
	splittedInfo := strings.Split(replacedInfo, " ")
	result := []*net.UDPAddr{}
	for _, value := range splittedInfo{
		if len(value) == 0{
			continue
		}
		splitedValue := strings.Split(value, ":")
		ip, portStr := splitedValue[0], splitedValue[1]
		port, err := strconv.Atoi(portStr)
		if err != nil{
			panic(err)
		}
		result = append(result, &net.UDPAddr{
			IP: net.ParseIP(ip),
			Port: port,
		})
	}
	return result
}

