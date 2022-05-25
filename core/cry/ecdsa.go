package cry

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
)

func Sign(key ecdsa.PrivateKey, msg []byte) ([]byte,error) {
	sig, err := ecdsa.SignASN1(rand.Reader, &key, msg)
	if err != nil {
		return nil, err
	}
	return sig, nil
}


func Verify(key ecdsa.PublicKey, msg []byte, sig []byte) bool {
	valid := ecdsa.VerifyASN1(&key, msg, sig)
	return valid
}

func B64SKTo(b64Sk string) (ecdsa.PrivateKey,error){
	b, err := B64toByte(b64Sk)
	if err != nil {
		return ecdsa.PrivateKey{
			PublicKey:ecdsa.PublicKey{Curve:elliptic.P256()},
		}, err
	}

	z := big.Int{}
	z.SetBytes(b)
	return ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve:elliptic.P256()},
		D: &z ,
	}, nil
}

func B64PKTo(b64Pk string) (ecdsa.PublicKey,error) {
	b, err := B64toByte(b64Pk)
	if err != nil {
		return ecdsa.PublicKey{Curve:elliptic.P256()},err
	}
	point := BytesToPoint(b)
	return ecdsa.PublicKey{
		X: point.X,
		Y: point.Y,
		Curve:elliptic.P256(),
	},nil
}