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

func (h *Header) BytesFromSeed(id []byte, rest []byte) []byte {
	bytes := make([]byte, 12)

	idInt := (uint16(id[0]) << 8) | uint16(id[1])
	opcode := uint16(rest[0]) & 0x78 >> 3
	rd := uint16(rest[0]) & 0x01

	binary.BigEndian.PutUint16(bytes[0:], idInt)
	binary.BigEndian.PutUint16(bytes[2:], uint16(h.QR)<<15|uint16(opcode)<<11|uint16(h.AA)<<10|uint16(h.TC)<<9|uint16(rd)<<8|uint16(h.RA)<<7|uint16(h.Z)<<4|uint16(h.RCODE))
	binary.BigEndian.PutUint16(bytes[4:], h.QDCOUNT)
	binary.BigEndian.PutUint16(bytes[6:], h.ANCOUNT)
	binary.BigEndian.PutUint16(bytes[8:], h.NSCOUNT)
	binary.BigEndian.PutUint16(bytes[10:], h.ARCOUNT)

	return bytes
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

type Question struct {
	NAME  []byte
	TYPE  uint16
	CLASS uint16
}

func (q *Question) Bytes() []byte {
	label := append(q.NAME, 0)
	bytes := make([]byte, len(label)+4)

	copy(bytes, label)
	binary.BigEndian.PutUint16(bytes[len(label):], q.TYPE)
	binary.BigEndian.PutUint16(bytes[len(label)+2:], q.CLASS)

	return bytes
}

type Answer struct {
	NAME     []byte
	TYPE     uint16
	CLASS    uint16
	TTL      uint32
	RDLENGTH uint16
	RDATA    []byte
}

func (a *Answer) Bytes() []byte {
	label := append(a.NAME, 0)
	bytes := make([]byte, len(label)+14)

	copy(bytes, label)

	binary.BigEndian.PutUint16(bytes[len(label):], a.TYPE)
	binary.BigEndian.PutUint16(bytes[len(label)+2:], a.CLASS)
	binary.BigEndian.PutUint32(bytes[len(label)+4:], a.TTL)
	binary.BigEndian.PutUint16(bytes[len(label)+8:], a.RDLENGTH)

	bytes = append(bytes, a.RDATA...)

	return bytes
}

type Authority struct{}

type Additional struct{}

type Message struct {
	Header      Header
	Questions   []Question
	Answers     []Answer
	Authorities []Authority
	Additionals []Additional
}

func (m *Message) Bytes() []byte {
	bytes := m.Header.Bytes()

	for _, question := range m.Questions {
		bytes = append(bytes, question.Bytes()...)
	}

	for _, answer := range m.Answers {
		bytes = append(bytes, answer.Bytes()...)
	}

	return bytes
}

func (m *Message) BytesFromSeed(id []byte, rest []byte) []byte {
	bytes := m.Header.BytesFromSeed(id, rest)

	for _, question := range m.Questions {
		bytes = append(bytes, question.Bytes()...)
	}

	for _, answer := range m.Answers {
		bytes = append(bytes, answer.Bytes()...)
	}

	return bytes
}

func domainNameListToHexString(domainName []string) []byte {
	var domainNameBytes []byte

	for _, part := range domainName {
		domainNameBytes = append(domainNameBytes, byte(len(part)))
		domainNameBytes = append(domainNameBytes, []byte(part)...)
	}

	return domainNameBytes
}
