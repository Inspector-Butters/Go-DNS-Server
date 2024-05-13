package main

import "encoding/binary"

type Header struct {
	ID      uint16
	QR      uint8
	OPCODE  uint8
	AA      uint8
	TC      uint8
	RD      uint8
	RA      uint8
	Z       uint8
	RCODE   uint8
	QDCOUNT uint16
	ANCOUNT uint16
	NSCOUNT uint16
	ARCOUNT uint16
}

func (h *Header) Bytes() []byte {
	bytes := make([]byte, 12)

	binary.BigEndian.PutUint16(bytes[0:], h.ID)
	binary.BigEndian.PutUint16(bytes[2:], uint16(h.QR)<<15|uint16(h.OPCODE)<<11|uint16(h.AA)<<10|uint16(h.TC)<<9|uint16(h.RD)<<8|uint16(h.RA)<<7|uint16(h.Z)<<4|uint16(h.RCODE))
	binary.BigEndian.PutUint16(bytes[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(bytes[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(bytes[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(bytes[10:], h.ARCOUNT)

	return bytes
}

type Question struct{}

type Answer struct{}

type Authority struct{}

type Additional struct{}

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Answer
	Authorities []Authority
	Additionals []Additional
}

