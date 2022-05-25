package cry

import (
	"fmt"
	"testing"
)

func TestAlg(t *testing.T) {
	sk,pk,_ := GenKey()
	fmt.Println(sk,pk)
	msg := []byte("Gophers, gophers, gophers everywhere!\n\r\n\nGophers, gophers, gophers everywhere!\n\n\nGophers, gophers, gophers everywhere!\n\n\n")
	key := []byte("12345678901234567890123456789012")
	ciphertext, err := Encrypt(msg, key)
	if err != nil {
		panic(err)
	}

	plaintext, err := Decrypt(ciphertext,key)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", plaintext)
}