package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	ascii "github.com/galsondor/go-ascii"
)

func getInfoFromFile()string{
	
	file, err := os.OpenFile("server.txt", os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil{
		panic(err)
	}
	buff := make([]byte, 4096)
	_, err = file.Read(buff)
	if err != nil{
		panic(err)
	}
	if len(buff) == 0{
		fmt.Println("The configuration file is empty!")
		os.Exit(0)
	}

	defer file.Close()
	return string(buff)
}

func cleanInfoFromFile(info string)string{
	var cleanedInfo string
	for _, symbol := range info{
		if ascii.IsGraph(byte(symbol)) || ascii.IsSpace(byte(symbol)){
			cleanedInfo += string(symbol)
		}
	}
	return cleanedInfo
}

func parseInfoFromFile(info string)[]*net.UDPAddr{
	replacedInfo := strings.Replace(info, "\n", " ", -1)
	cleanedInfo := strings.Split(cleanInfoFromFile(replacedInfo), " ")
	result := []*net.UDPAddr{}
	for _, value := range cleanedInfo{
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

func createConnections(parsedInfo []*net.UDPAddr)[]net.Conn{
	result := []net.Conn{}
	for _, addr := range parsedInfo{
		
		conn, err := net.Dial("udp", addr.String())
		if err != nil{
			panic(err)
		}
		result = append(result, conn)
	}
	return result
}

func checkServersToReady(subServers []net.Conn)[]net.Conn{
	result := []net.Conn{}
	for _, server := range subServers{
		server.SetReadDeadline(time.Now().Add(1 * time.Second))
		server.Write([]byte("OK///"))
		buff := make([]byte, 1)
		server.Read(buff)
		if string(buff) == "1"{
			result = append(result, server)
		}
	}
	return result
}

func formatMessage(availableServers []net.Conn)string{
	result := []string{}
	for _, value := range availableServers{
		result = append(result, value.RemoteAddr().String())
	}
	return strings.Join(result, " ")
}

func getSubServers()[]net.Conn{
	freshInfo := getInfoFromFile()
	parsedInfo := parseInfoFromFile(freshInfo)
	createdConnections := createConnections(parsedInfo)
	return createdConnections
}

func main(){
	envaddr := os.Getenv("MAINSERVER_ADDR")
	if len(envaddr) == 0{
		fmt.Println("MAINSERVER_ADDR is not written!")
		os.Exit(0)
	}
	splittedAddr := strings.Split(envaddr, ":")
	ip, portStr := splittedAddr[0], splittedAddr[1]
	port, err := strconv.Atoi(portStr)
	if err != nil{
		panic(err)
	}

	connection, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP: net.ParseIP(ip),
	})
	if err != nil{
		panic(err)
	}
	subServers := getSubServers()
	for{
		buff := make([]byte, 4096)
		_, addr, err := connection.ReadFromUDP(buff)
		if err != nil{
			panic(err)
		}
		switch{
		case strings.Contains(string(buff), "CheckServers"):
			message := formatMessage(checkServersToReady(subServers))
			connection.WriteTo([]byte(message), addr)
		default:
			connection.WriteTo([]byte("Commands is not avaialable!\n"), addr)
		}
	}
}