package cry

import (
	"encoding/base64"
)

func ByteTob64(input []byte) string {
	return base64.StdEncoding.EncodeToString(input)
}

func B64toByte(b64 string) ([]byte,error){
	data,err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil,err
	}
	return data,nil
}