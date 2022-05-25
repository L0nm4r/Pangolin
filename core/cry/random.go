package cry

import cryptorand "crypto/rand"

func GetNonce() (uint32,error){
	nonce := make([]byte,4)
	_, err := cryptorand.Read(nonce)
	if err != nil {
		return 0, err
	}
	return BytesToUInt32(nonce),nil
}

func GetRandomBytes() ([]byte,error){
	nonce := make([]byte,8)
	_, err := cryptorand.Read(nonce)
	if err != nil {
		return nil, err
	}
	return nonce,nil
}