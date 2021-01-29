package main

import (
	"MainServer/Messageparser"
	"MainServer/Networking"
	"bytes"
	"io"
	"net"
	"os"
	"strings"
	"time"
)


func readyServers(sservers []net.Conn)[]*Messageparser.Message{
	var result []*Messageparser.Message
	result = append(result, new(Messageparser.Message))
	
	for _, server := range sservers{
		server.SetWriteDeadline(time.Now().Add(1 * time.Second))
		server.SetReadDeadline(time.Now().Add(1 * time.Second))

		parser := Messageparser.Parser(new(Messageparser.Message))
		b := parser.Parse(Messageparser.NewMessage())
		for{
			_, err := server.Write(b)
			
			if os.IsTimeout(err){
				continue
			}
		
			var buff bytes.Buffer
			_, err = io.Copy(&buff, server)

			if os.IsTimeout(err) && buff.Len() == 0{
				continue
			}

			if buff.Len() != 0{
				message, err := parser.Unparse(buff.Bytes())
				if err != nil{
					continue
				}
				if message[0].Error == "200"{
					result[0].Body = append(result[0].Body, server.RemoteAddr().String())
				}
				break
			}
		}
		server.Close()
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

func main(){

	sp := Messageparser.StartParser(new(Messageparser.SP))
	inf := sp.GetStart()
	ip, port := sp.ParseStart(inf)

	connection, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: port,
		IP: net.ParseIP(ip),
	})
	if err != nil{
		panic(err)
	}

	sserversparser := Messageparser.FileParser(new(Messageparser.FP))
	sservers := sserversparser.ParseFile(sserversparser.GetFile())

	for{
		notCleaned := make([]byte, 4096)
		_, addr, err := connection.ReadFromUDP(notCleaned)
		var buff []byte
		for _, value := range notCleaned{
			if value != 0{
				buff = append(buff, value)
			}
		}
		if len(buff) != 0{
			if err != nil{
				return
			}
			parser := Messageparser.Parser(new(Messageparser.Message))
			message, err := parser.Unparse(buff)
			if err != nil{
				connection.WriteTo([]byte("Commands is not avaialable!\n"), addr)
				continue
			}
			switch{
			case message[0].Type == "CheckServers":
				go func(){
					readyServers := readyServers(Networking.CreateConnections(sservers))
					readyServers[0].Index = message[0].Index
					b := parser.Parse(readyServers)
					connection.WriteTo(b, addr)
				}()
			default:
				connection.WriteTo([]byte("Commands is not avaialable!\n"), addr)
			}
		}
	}
}
