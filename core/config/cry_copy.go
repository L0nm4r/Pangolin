package config

import (
	"crypto/sha1"
	"encoding/hex"
)

func HashCodeStr(s string) []byte {
	data := []byte(s)
	m := sha1.New()
	m.Write(data)
	return m.Sum([]byte{})
}

func BytesToString(b []byte) string {
	return hex.EncodeToString(b)
}