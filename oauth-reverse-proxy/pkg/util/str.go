package util

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
)

var baseRune = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var baseRuneLength = uint32(len(baseRune))

func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func RandomString(n int) string {
	var bb bytes.Buffer
	bb.Grow(n)
	for i := 0; i < n; i++ {
		bb.WriteRune(baseRune[binary.BigEndian.Uint32(Bytes(4))%baseRuneLength])
	}
	return bb.String()
}
