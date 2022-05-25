package cry

import (
	cryptorand "crypto/rand"
	"fmt"
	alg "golang.org/x/crypto/chacha20poly1305"
)

func example() {
	// key should be randomly generated or derived from a function like Argon2.
	KeySize := 32
	key := make([]byte, KeySize)
	if _, err := cryptorand.Read(key); err != nil {
		panic(err)
	}

	fmt.Println(ByteTob64(key))
	aead, err := alg.NewX(key)
	if err != nil {
		panic(err)
	}
	// Encryption.
	var encryptedMsg []byte
	{
		msg := []byte("Gophers, gophers, gophers everywhere!")

		// Select a random nonce, and leave capacity for the ciphertext.
		nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(msg)+aead.Overhead())
		if _, err := cryptorand.Read(nonce); err != nil {
			panic(err)
		}

		// Encrypt the message and append the ciphertext to the nonce.
		encryptedMsg = aead.Seal(nonce, nonce, msg, nil)
	}

	// Decryption.
	{
		if len(encryptedMsg) < aead.NonceSize() {
			panic("ciphertext too short")
		}

		// Split nonce and ciphertext.
		nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

		// Decrypt the message and check it wasn't tampered with.
		plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", plaintext)
	}
}