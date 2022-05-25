package cry

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"github.com/aead/ecdh"
)

// ECDH key-exchange using the curve P256.
// https://pkg.go.dev/github.com/aead/ecdh#section-readme
func GenKey() ([]byte ,[]byte,error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil,nil, err
	}
	//fmt.Println("GenKey: ","\nprivate key:",privateKey,"\npublic key:",privateKey.PublicKey)
	return privateKey.D.Bytes(), PointToBytes(privateKey.PublicKey),nil
}

func CalculateSecret(privateAlice ecdsa.PrivateKey,publicBob ecdsa.PublicKey) ([]byte,error) {
	point := ecdh.Point {
		X: publicBob.X,
		Y: publicBob.Y,
	}
	p256 := ecdh.Generic(elliptic.P256())
	// Alice
	if err := p256.Check(point); err != nil {
		//fmt.Printf("Bob's public key is not on the curve: %s\n", err)
		return nil,err
	}
	secret := p256.ComputeSecret(privateAlice.D.Bytes(), point)
	return secret,nil
}