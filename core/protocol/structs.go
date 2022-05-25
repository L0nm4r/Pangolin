package protocol

import (
	"Pangolin/core/config"
	cry2 "Pangolin/core/cry"
	"bytes"
	"crypto/ecdsa"
	"errors"
	"fmt"
)

type HandShakeRequest struct {
	MsgType byte		// 1 byte
	ClientID []byte		// 20 Size Sha1 Of Client' PublicKey
	TimeStamp uint64	//
	Nonce uint32		//
	SPk []byte			// 68 byte
	HashCode []byte 	// 20 size // Sha1(MsgType || ClientID || TimeStamp || Nonce || SPk || Sha1(Server's PublicKey || Client's PublicKey))
	Signature []byte		// size ?? ECDSA_SIGN(HASHCODE)
}

func (r *HandShakeRequest) FillHash() {
	hash1 := cry2.HashCodeStr(config.ServerPK + config.ClientPK)
	b := []byte{r.MsgType}
	b = append(b, r.ClientID...)
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, hash1...)
	r.HashCode = cry2.HashCode(b)
}

func (r *HandShakeRequest) Sign(sk string) error {
	key, err := cry2.B64SKTo(sk)
	if err != nil {
		return errors.New("convert b64 sk to bytes error")
	}
	signature, err := cry2.Sign(key,r.HashCode)
	if err != nil {
		return errors.New("sign error")
	}
	r.Signature = signature
	return nil
}

func (r *HandShakeRequest) ToBytes() []byte {
	b := []byte{r.MsgType}
	b = append(b, r.ClientID...)
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, r.HashCode...)
	b = append(b, r.Signature...)

	return b
}

func (r *HandShakeRequest) From(buf []byte) error {
	r.MsgType = buf[0]
	r.ClientID = buf[1:21]
	r.TimeStamp = cry2.BytesToUInt64(buf[21:29])
	r.Nonce = cry2.BytesToUInt32(buf[29:33])
	r.SPk = buf[33:101]
	r.HashCode = buf[101:121]
	r.Signature = buf[121:]
	return nil
}

func (r *HandShakeRequest) VerifyTimestamp() (bool,error) {
	ok := cry2.CheckTimeStamp(int64(r.TimeStamp))
	return ok,nil
}

func (r *HandShakeRequest) VerifyHashCode(clientPK string) (bool, error) {
	hash1 := cry2.HashCodeStr(config.ServerPK + clientPK)
	b := []byte{r.MsgType}
	b = append(b, r.ClientID...)
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, hash1...)
	hashcode := cry2.HashCode(b)
	if bytes.Equal(hashcode[:20],r.HashCode[:20]) {
		return true, nil
	}
	return false,nil
}

func (r *HandShakeRequest) VerifySignature(key ecdsa.PublicKey) (bool, error) {
	ok := cry2.Verify(key,r.HashCode,r.Signature)
	return ok,nil
}

type HandShakeResponse struct {
	MsgType byte		// 2
	TimeStamp uint64		//
	Nonce uint32			// request's nonce + 1
	SPk []byte		// 68 size Server's Session Public Key
	HashCode []byte 	//  20 size Hash(MsgType || TimeStamp || Nonce || SPk || Hash(Server's PublicKey || Client's PublicKey))
	Signature []byte  	// 71 or 70 size Client -> Sign(Hash(HandShakeResponseEnc + MsgType + TimeStamp)))
}

func (r *HandShakeResponse) From(buf []byte) error {
	r.MsgType = buf[0]
	r.TimeStamp = cry2.BytesToUInt64(buf[1:9])
	r.Nonce = cry2.BytesToUInt32(buf[9:13])
	r.SPk = buf[13:81]
	r.HashCode = buf[81:101]
	r.Signature = buf[101:]
	return nil
}

func (r *HandShakeResponse) VerifyTimeStamp() (bool,error) {
	ok := cry2.CheckTimeStamp(int64(r.TimeStamp))
	return ok,nil
	//return true,nil // TODO : debug
}

func (r *HandShakeResponse) VerifyNonce(nonce uint32) (bool, error) {
	return r.Nonce == nonce + 1,nil
}

func (r *HandShakeResponse) VerifyHash() (bool, error) {
	hash1 := cry2.HashCodeStr(config.ServerPK + config.ClientPK)
	b := []byte{r.MsgType}
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, hash1...)
	hashcode := cry2.HashCode(b)
	return bytes.Equal(r.HashCode,hashcode),nil
}

func (r *HandShakeResponse) VerifySignature(key ecdsa.PublicKey) (bool, error) {
	ok := cry2.Verify(key, r.HashCode, r.Signature)
	return ok,nil
}

func (r *HandShakeResponse) ToBytes() []byte {
	b := []byte{r.MsgType}
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, r.HashCode...)
	b = append(b, r.Signature...)
	return b
}

func (r *HandShakeResponse) FillHash(clientPK string) error {
	hash1 := cry2.HashCodeStr(config.ServerPK + clientPK)
	b := []byte{r.MsgType}
	b = append(b, cry2.UInt64ToBytes(r.TimeStamp)...)
	b = append(b, cry2.UInt32ToBytes(r.Nonce)...)
	b = append(b, r.SPk...)
	b = append(b, hash1...)
	r.HashCode = cry2.HashCode(b)

	return nil
}

func (r *HandShakeResponse) Sign(sk string) error {
	key, err := cry2.B64SKTo(sk)
	if err != nil {
		return err
	}
	r.Signature, err = cry2.Sign(key,r.HashCode)
	return err
}

type DataTransport struct {
	MsgType byte			// 3
	Counter byte			// Counter ++
	Timestamp uint64		// 8 byte
	HashCode []byte			// 16 byte HMAC(SessionKey,(MsgType || Counter || Data))
	Data []byte
}

func (d *DataTransport) ToBytes() []byte {
	b := []byte{d.MsgType,d.Counter}
	b = append(b, cry2.UInt64ToBytes(d.Timestamp)...)
	b = append(b, d.HashCode...)
	b = append(b, d.Data...)
	return b
}

func (d *DataTransport) FromBytes(buf []byte) error {
	if len(buf) < 26 {
		return errors.New("bad data length")
	}
	d.MsgType = buf[0]
	d.Counter = buf[1]

	d.Timestamp = cry2.BytesToUInt64(buf[2:10])

	d.HashCode = buf[10:26]
	d.Data = buf[26:]
	return nil
}

func (d *DataTransport) CalcHash(sessionKey []byte) error {
	b := []byte{d.MsgType,d.Counter}
	b = append(b, cry2.UInt64ToBytes(d.Timestamp)...)
	b = append(b, d.Data...)

	d.HashCode = cry2.HmacMd5(sessionKey, b)
	//fmt.Printf("send hashcode: %s\n", d.HashCode)
	return nil
}

func (d *DataTransport) Verify(expectCounter byte,sessionKey []byte) (bool,error){
	if !cry2.CheckTimeStamp(int64(d.Timestamp)) ||expectCounter != d.Counter {
		return false,errors.New(fmt.Sprintf("time or counter verify error"))
	}

	b := []byte{d.MsgType,d.Counter}
	b = append(b, cry2.UInt64ToBytes(d.Timestamp)...)
	b = append(b, d.Data...)

	hash := cry2.HmacMd5(sessionKey, b)
	if !bytes.Equal(hash,d.HashCode) {
		return false,errors.New(fmt.Sprintf("hash calc error"))
	}

	return true,nil
}