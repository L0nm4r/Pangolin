package protocol

import (
	"Pangolin/core/config"
	cry2 "Pangolin/core/cry"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"math/big"
	"net"
	"sync"
	"time"
)


// GetConnect : HandShake Request and response
func GetConnect(addr string) (Tunnel, error) { // session Key, connection, error
	// Gen Session Key
	sk,pk,err := cry2.GenKey()
	if err != nil {
		return Tunnel{}, errors.New("gen session key error")
	}

	// HandShake
	server, err := net.DialTimeout("tcp",addr, 3*time.Second)
	tunnel := Tunnel{
		SessionKey: nil,
		Counter:    0,
		Conn:       server,
		Mutex:      sync.Mutex{},
		Buffer: map[uint64][]byte{},
	}
	if err != nil {
		return tunnel,err
	}
	ts := cry2.GenTimeStamp()
	nonce,err := cry2.GetNonce()
	if err != nil {
		return tunnel,errors.New("get nonce error")
	}

	var request = HandShakeRequest{
		MsgType:   1,
		ClientID:  cry2.HashCodeStr(config.ClientPK),
		TimeStamp: uint64(ts),
		Nonce:     nonce,
		SPk:       pk,
	}

	request.FillHash()
	err = request.Sign(config.ClientSK)
	if err != nil {
		return tunnel, err
	}

	_, err = server.Write(request.ToBytes())
	if err != nil {
		return tunnel, err
	}

	buf := make([]byte, 10240) // response max size ?
	n := 0
	for true {
		n, err = server.Read(buf)
		if err != nil {
			return tunnel, err
		}
		if n != 0 {
			break
		}
	}
	//  handle response
	response := HandShakeResponse{}
	err = response.From(buf[:n])
	if err != nil {
		return tunnel,err
	}

	if int(response.MsgType) != 2 {
		return tunnel,errors.New("msg type error")
	}

	ok,err := response.VerifyTimeStamp()
	if err != nil || !ok {
		return tunnel,errors.New("time stamp verify error")
	}

	ok,err = response.VerifyNonce(nonce)
	if err != nil || !ok {
		return tunnel,errors.New("nonce verify error")
	}

	ok,err = response.VerifyHash()
	if err != nil || !ok {
		return tunnel,errors.New("hash verify error")
	}

	ServerPublicKey,err := cry2.B64PKTo(config.ServerPK)
	ok,err = response.VerifySignature(ServerPublicKey)
	if err != nil || !ok {
		return tunnel,errors.New("signature verify error")
	}

	// session Key
	AlicePK := cry2.BytesToPoint(response.SPk)
	tsk := big.Int{}
	tsk.SetBytes(sk)
	BobSessionPK := cry2.BytesToPoint(pk)
	if err != nil {
		return tunnel, err
	}
	BobSessionSK := ecdsa.PrivateKey{
		D: &tsk,
		PublicKey:BobSessionPK,
	}

	sessionKey, err := cry2.CalculateSecret(BobSessionSK,AlicePK)
	if err != nil {
		return tunnel, err
	}
	tunnel.SessionKey = sessionKey

	return tunnel, nil
}

// 函数调用者:defer client.close()
func HandShakeHander(client net.Conn) ([]byte, error) { // session Key, connection, error
	var buf = make([]byte,10240) // maxsize?
	n,err := client.Read(buf)
	if err != nil {
		return nil,err
	}
	var request = HandShakeRequest{}
	err = request.From(buf[:n])
	if err != nil {
		return nil,err
	}
	// close ??
	if request.MsgType != 1 {
		return nil,errors.New("msg type error")
	}

	hashcode := BytesToString(request.ClientID)

	clientPK := config.Clients[hashcode] // error

	clientPublicKey,err := cry2.B64PKTo(clientPK)
	if err != nil {
		return nil,errors.New("client decode error")
	}

	ok,err := request.VerifyTimestamp()
	if err != nil || !ok {
		return nil,errors.New("timestamp verify error")
	}

	ok,err = request.VerifyHashCode(clientPK)
	if err != nil || !ok {
		return nil,errors.New("hashcode verify error")
	}

	ok,err = request.VerifySignature(clientPublicKey)
	if err != nil || !ok {
		return nil,errors.New("signature verify error")
	}

	// gen session key
	sk,pk,err := cry2.GenKey()
	if err != nil {
		return nil, errors.New("gen session key error")
	}
	// session Key
	AlicePK := cry2.BytesToPoint(request.SPk)
	tsk := big.Int{}
	tsk.SetBytes(sk)
	BobSessionPK := cry2.BytesToPoint(pk)
	if err != nil {
		return nil,err
	}
	BobSessionSK := ecdsa.PrivateKey{
		D: &tsk,
		PublicKey:BobSessionPK,
	}

	sessionKey, err := cry2.CalculateSecret(BobSessionSK,AlicePK)
	if err != nil {
		return nil, err
	}

	response := HandShakeResponse{
		MsgType:   2,
		TimeStamp: uint64(cry2.GenTimeStamp()),
		Nonce:     request.Nonce + 1,
		SPk:       pk,
	}

	err = response.FillHash(clientPK)
	if err != nil {
		return nil,err
	}

	err = response.Sign(config.ServerSK)
	if err != nil {
		return nil,err
	}

	_,err = client.Write(response.ToBytes())
	if err != nil {
		return nil, err
	}
	return sessionKey, nil
}

func BytesToString(b []byte) string {
	return hex.EncodeToString(b)
}