package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)
		//print the received data
		fmt.Println(receivedData)

		response := Message{
			Header: Header{
				ID:      1234,
				QR:      1,
				OPCODE:  0,
				AA:      0,
				TC:      0,
				RD:      0,
				RA:      0,
				Z:       0,
				RCODE:   0,
				QDCOUNT: 1,
				ANCOUNT: 0,
				NSCOUNT: 0,
				ARCOUNT: 0,
			},
			Questions: []Question{
				{
					NAME:  []byte("\x0ccodecrafters\x02io"),
					TYPE:  1,
					CLASS: 1,
				},
			},
		}

		_, err = udpConn.WriteToUDP(response.Bytes(), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
