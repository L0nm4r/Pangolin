package cry

import (
	cryptorand "crypto/rand"
	alg "golang.org/x/crypto/chacha20poly1305"
)

// chacha20-ietf-poly1305
// https://pkg.go.dev/golang.org/x/crypto/chacha20poly1305

func Encrypt(msg []byte, sessionKey []byte) ([]byte,error) {
	key,err := alg.NewX(sessionKey)
	if err != nil {
		return nil, err
	}
	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, key.NonceSize(), key.NonceSize()+len(msg)+key.Overhead())
	if _, err := cryptorand.Read(nonce); err != nil {
		return nil, err
	}
	// Encrypt the message and append the ciphertext to the nonce.
	encryptedMsg := key.Seal(nonce, nonce, msg, nil)
	return encryptedMsg,nil
}

func Decrypt(encryptedMsg []byte, sessionKey []byte) ([]byte,error) {
	//  处理异常的函数
	key,err := alg.NewX(sessionKey)
	if err != nil {
		return nil, err
	}
	nonce, ciphertext := encryptedMsg[:key.NonceSize()], encryptedMsg[key.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := key.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	//fmt.Printf("%s\n", plaintext)
	return plaintext,nil
}