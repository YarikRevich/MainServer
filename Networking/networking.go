package Networking

import(
	"net"
)

func CreateConnections(parsedInfo []*net.UDPAddr)[]net.Conn{
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