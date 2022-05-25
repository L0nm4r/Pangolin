package cry

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
)

func HashCode(data []byte) []byte {
	m := sha1.New()
	m.Write(data)
	return m.Sum([]byte{})
}

func HashCodeStr(s string) []byte {
	data := []byte(s)
	m := sha1.New()
	m.Write(data)
	return m.Sum([]byte{})
}

func HmacMd5(key, data []byte) []byte {
	h := hmac.New(md5.New, key)
	h.Write(data)
	return h.Sum([]byte("")) // 128位, 16字节
}

func HmacMd5Str(key, data string) string {
	h := hmac.New(md5.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum([]byte("")))
}