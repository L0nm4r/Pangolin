package cry

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestEcdsa(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}

	fmt.Printf("public Key: %s:%s\nprivate Key: %s\n",privateKey.PublicKey.X,privateKey.PublicKey.Y,privateKey.D)

	msg := "hello, world"
	hash := sha256.Sum256([]byte(msg))

	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	fmt.Printf("signature: %x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	fmt.Println("signature verified:", valid)
}

func TestBytesToPoint(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	b := PointToBytes(privateKey.PublicKey)
	k := BytesToPoint(b)

	fmt.Println(k)
	fmt.Println(privateKey.PublicKey)
}

func TestKeyGender(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	pk := PointToBytes(privateKey.PublicKey)
	sk := privateKey.D.Bytes()

	fmt.Println("publicKey:", ByteTob64(pk))
	fmt.Println("privateKey:", ByteTob64(sk))
}