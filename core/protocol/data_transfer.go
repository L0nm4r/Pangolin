package protocol

import (
	cry2 "Pangolin/core/cry"
	"errors"
	"fmt"
)

func PackData(raw []byte, SessionKey []byte, counter byte) ([]byte,error){
	cipherData, err := cry2.Encrypt(raw,SessionKey)
	if err != nil {
		return nil,err
	}

	transport := DataTransport{
		MsgType:   3,
		Counter:   counter,
		Timestamp: uint64(cry2.GenTimeStamp()),
		Data:      cipherData,
	}
	err = transport.CalcHash(SessionKey)
	if err != nil {
		return nil, err
	}
	return transport.ToBytes(), err
}

func UnPackData(raw []byte, SessionKey []byte, expectedCounter byte) ([]byte,error) {
	data := DataTransport{}
	err := data.FromBytes(raw)
	if err != nil {
		return nil, err
	}
	ok ,err := data.Verify(expectedCounter,SessionKey)
	if !ok {
		return nil,errors.New(fmt.Sprintf("bad data format %s", err))
	}
	plainText, err := cry2.Decrypt(data.Data, SessionKey)
	if err != nil {
		return nil, err
	}
	return plainText,nil
}
