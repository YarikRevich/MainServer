package Messageparser

import(
	"os"
	"fmt"
	"strings"
	"strconv"
)

type StartParser interface{
	GetStart()string
	ParseStart(string)(string, int)
}

type SP struct{}

func(s *SP) GetStart()string{
	envaddr := os.Getenv("MAINSERVER_ADDR")
	if len(envaddr) == 0{
		fmt.Println("MAINSERVER_ADDR is not written!")
		os.Exit(0)
	}
	return envaddr
}

func(s *SP)ParseStart(info string)(string, int){
	splittedAddr := strings.Split(info, ":")
	ip, portStr := splittedAddr[0], splittedAddr[1]
	port, err := strconv.Atoi(portStr)
	if err != nil{
		panic(err)
	}
	return ip, port
}
