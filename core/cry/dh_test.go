package cry

import (
	"bytes"
	"crypto"
	"crypto/elliptic"
	"crypto/rand"
	"fmt"
	"github.com/aead/ecdh"
	"github.com/vmihailenco/msgpack"
	"math/big"
	"testing"
)

func TestDH(t *testing.T) {
	p256 := ecdh.Generic(elliptic.P256())

	// Alice
	privateAlice, publicAlice, err := p256.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate Alice's private/public key pair: %s\n", err)
	}

	pA, err := msgpack.Marshal(publicAlice)
	if err != nil {
		panic(err)
	}

	// privateAlice []uint8
	// publicAlice ecdh.Point
	// Bob
	privateBob, publicBob, err := p256.GenerateKey(rand.Reader)
	if err != nil {
		fmt.Printf("Failed to generate Bob's private/public key pair: %s\n", err)
	}

	pB, err := msgpack.Marshal(publicBob)
	if err != nil {
		panic(err)
	}

	// Alice
	var pBB = ecdh.Point{}
	err = msgpack.Unmarshal(pB, &pBB)

	if err != nil {
		panic(err)
	}
	if err := p256.Check(crypto.PublicKey(pBB)); err != nil {
		fmt.Printf("Bob's public key is not on the curve: %s\n", err)
	}
	secretAlice := p256.ComputeSecret(privateAlice, pBB)

	// Bob
	var pAA = ecdh.Point{}
	err = msgpack.Unmarshal(pA, &pAA)
	if err != nil {
		panic(err)
	}
	if err := p256.Check(crypto.PublicKey(pAA)); err != nil {
		fmt.Printf("Alice's public key is not on the curve: %s\n", err)
	}
	secretBob := p256.ComputeSecret(privateBob, publicAlice)

	fmt.Printf("secretAlice: %s\nsecretBob: %s\n",secretAlice, secretBob)
	if !bytes.Equal(secretAlice, secretBob) {
		fmt.Printf("key exchange failed - secret X coordinates not equal\n")
	}
}

func PublicKey2Bytes(key ecdh.Point) (int, int, []byte) {
	var buf []byte
	var xLen = big.NewInt(int64(len(key.X.Bytes())))
	var yLen = big.NewInt(int64(len(key.Y.Bytes())))
	buf = append(buf,xLen.Bytes()...)
	buf = append(buf,yLen.Bytes()...)
	buf = append(buf,key.X.Bytes()...)
	buf = append(buf,key.Y.Bytes()...)
	return len(key.X.Bytes()), len(key.Y.Bytes()), buf
}