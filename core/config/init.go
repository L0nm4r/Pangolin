package config

import (
	"encoding/json"
)

type PK struct {
	PublicKey string `json:"publicKey"`
}

type Config struct {
	Server PK `json:"server"`
	Client PK `json:"client"`
	ClientPrivateKey string `json:"clientPrivateKey"`
	ServerPrivateKey string `json:"serverPrivateKey"`
	Clients []PK `json:"clients"`
}

func InitConfig(bytes []byte) error {
	var C = Config{
		Server:           PK{},
		ClientPrivateKey: "",
		ServerPrivateKey: "",
		Clients:          nil,
	}
	err := json.Unmarshal([]byte(bytes), &C)
	if err != nil {
		return err
	}

	ClientPK = C.Client.PublicKey
	ServerPK = C.Server.PublicKey
	ClientSK = C.ClientPrivateKey
	ServerSK = C.ServerPrivateKey

	for _,k := range C.Clients {
		Clients[BytesToString(HashCodeStr(k.PublicKey))] =  k.PublicKey
	}
	return nil
}