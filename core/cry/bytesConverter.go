package cry

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/binary"
	"math/big"
)

func PointToBytes(point ecdsa.PublicKey) []byte {
	var b []byte
	xLen := uint32(len(point.X.Bytes()))
	b = append(b, UInt32ToBytes(xLen)...)
	b = append(b, point.X.Bytes()...)
	b = append(b, point.Y.Bytes()...)
	return b
}

func BytesToPoint(b []byte) ecdsa.PublicKey {
	publicKey := ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     &big.Int{},
		Y:     &big.Int{},
	}
	xLen := int(BytesToUInt32(b[:4]))
	publicKey.X.SetBytes(b[4:4+xLen])
	publicKey.Y.SetBytes(b[xLen+4:len(b)])

	return publicKey
}

func UInt32ToBytes(n uint32) []byte {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return b
}

func BytesToUInt32(b []byte) uint32 {
	n := binary.BigEndian.Uint32(b)
	return n
}

func UInt64ToBytes(n uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return b
}

func BytesToUInt64(b []byte) uint64 {
	n := binary.BigEndian.Uint64(b)
	return n
}