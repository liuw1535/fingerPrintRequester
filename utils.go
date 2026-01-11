package main

import (
	"crypto/rand"
	"encoding/binary"
	"strconv"
	"strings"
)

func generateRandomPSKIdentity() []byte {
	identity := make([]byte, 138)
	rand.Read(identity)
	return identity
}

func generateRandomObfuscatedTicketAge() uint32 {
	var age uint32
	binary.Read(rand.Reader, binary.BigEndian, &age)
	return age
}

func generateRandomBinder() []byte {
	binder := make([]byte, 33)
	rand.Read(binder)
	return binder
}

func generateGREASEValue() uint16 {
	greaseValues := []uint16{
		0x0a0a, 0x1a1a, 0x2a2a, 0x3a3a,
		0x4a4a, 0x5a5a, 0x6a6a, 0x7a7a,
		0x8a8a, 0x9a9a, 0xaaaa, 0xbaba,
		0xcaca, 0xdada, 0xeaea, 0xfafa,
	}
	var idx byte
	rand.Read([]byte{idx})
	return greaseValues[int(idx)%len(greaseValues)]
}

func parseHex(s string) uint16 {
	s = strings.TrimPrefix(s, "0x")
	val, _ := strconv.ParseUint(s, 16, 16)
	return uint16(val)
}
