package main

import (
	"fmt"
	"net"
	"strings"
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

		receivedData := buf[:size]
		fmt.Printf("Received %d bytes from %s\n", size, source)
		// fmt.Printf("Data: %v\n", receivedData)

		questionPart := receivedData[12:]
		var nameParts []string

		for i := 0; i < len(questionPart); i++ {
			length := int(questionPart[i])
			if length == 0 {
				break
			}
			i++
			nameParts = append(nameParts, string(questionPart[i:i+length]))
			i += length
		}
		domainName := []byte(strings.Join(nameParts, "."))

		idInt := (uint16(receivedData[0]) << 8) | uint16(receivedData[1])
		opcode := uint16(receivedData[2]) & 0x78 >> 3
		rd := uint16(receivedData[2]) & 0x01

		fmt.Println("id", idInt)
		fmt.Println("opcode", opcode)
		fmt.Println("domain name", domainName, nameParts)

		response := Message{
			Header: Header{
				ID:      idInt,
				QR:      1,
				OPCODE:  uint8(opcode),
				AA:      0,
				TC:      0,
				RD:      uint8(rd),
				RA:      0,
				Z:       0,
				RCODE:   4,
				QDCOUNT: 1,
				ANCOUNT: 1,
				NSCOUNT: 0,
				ARCOUNT: 0,
			},
			Questions: []Question{
				{
					NAME:  domainName,
					TYPE:  1,
					CLASS: 1,
				},
			},
			Answers: []Answer{
				{
					NAME:     domainName,
					TYPE:     1,
					CLASS:    1,
					TTL:      60,
					RDLENGTH: 4,
					RDATA:    []byte("\x08\x08\x08\x08"),
				},
			},
		}

		_, err = udpConn.WriteToUDP(response.BytesFromSeed(receivedData[0:2], receivedData[2:4]), source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}
