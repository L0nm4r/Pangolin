package protocol

import (
	cry2 "Pangolin/core/cry"
	"bytes"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
)

var DictPacket = []byte{0,0}
var ForwardSlice = []byte{1,0}
var BackSlice = []byte{0,1}

type Tunnel struct {
	SessionKey []byte
	Counter byte
	Conn  net.Conn
	Mutex sync.Mutex	// no use
	Buffer map[uint64][]byte	// no use
	LocalCounter byte
}

func (t *Tunnel) Read() (int,[]byte,error) {
	var buf = make([]byte,1<<16) // max size
	n,err := t.Conn.Read(buf)
	if err != nil || n == 0 {
		return 0,nil,err
	}
	data, err := UnPackData(buf[:n], t.SessionKey,t.Counter)

	// hash calc error 太多次了
	if err != nil {
		if strings.HasPrefix(fmt.Sprintf("%s",err), "bad data format") {
			return 0,nil,nil
		} else {
			return 0,nil,err
		}
	}
	t.Counter = t.Counter + 1

	dataType := data[:2]

	if bytes.Equal(dataType, DictPacket) {
		return len(data)-2,data[2:],nil
	}

	seq := cry2.BytesToUInt64(data[2:10])

	if bytes.Equal(dataType, ForwardSlice) {
		t.Buffer[seq] = data[10:]
		return 0,nil,nil
	} else if bytes.Equal(dataType, BackSlice) {
		return len(t.Buffer[seq])+len(data[10:]), append(t.Buffer[seq],data[10:]...), nil
	} else {
		return 0,nil,errors.New("error data type")
	}
}

//func (t *Tunnel) Read() (int,[]byte,error) {
//	var buf = make([]byte,1<<16) // max size
//	n,err := t.Conn.Read(buf)
//	if err != nil || n == 0 {
//		return 0,nil,err
//	}
//	t.Counter = t.Counter + 1
//	data, err := UnPackData(buf[:n], t.SessionKey,t.Counter)
//	fmt.Printf("read data size: %d, unpacked data: %d\n", n,len(data))
//	if err != nil {
//		t.Counter = t.Counter - 1
//		return 0,nil,err
//	}
//
//	return len(data),data,nil
//}

func (t *Tunnel) Write(msg[] byte) error {
	if len(msg) > 65000 {
		seq, err := cry2.GetRandomBytes()
		if err != nil {
			return err
		}

		idx := len(msg)/2
		forwardPadding := append(ForwardSlice,seq...)
		backPadding := append(BackSlice, seq...)

		// forward data
		packedData, err := PackData(append(forwardPadding,msg[:idx]...),t.SessionKey,t.LocalCounter)
		if err != nil {
			return err
		}
		_, err = t.Conn.Write(packedData)
		if err != nil {
			return err
		}
		t.LocalCounter = t.LocalCounter + 1

		// back data
		packedData, err = PackData(append(backPadding,msg[idx:]...),t.SessionKey,t.LocalCounter)
		if err != nil {
			return err
		}
		_, err = t.Conn.Write(packedData)
		if err != nil {
			return err
		}
		t.LocalCounter = t.LocalCounter + 1

	} else {
		packedData, err := PackData(append(DictPacket,msg...),t.SessionKey,t.LocalCounter)
		if err != nil {
			return err
		}
		_, err = t.Conn.Write(packedData)
		if err != nil {
			return err
		}
		t.LocalCounter = t.LocalCounter + 1
	}

	return nil
}

//func (t *Tunnel) Write(msg[] byte) error {
//	t.Counter = t.Counter + 1
//	packedData, err := PackData(msg,t.SessionKey,t.Counter)
//	fmt.Printf("write data size: %d, packeted data: %d\n", len(msg), len(packedData))
//	if err != nil {
//		t.Counter = t.Counter - 1
//		return err
//	}
//	_, err = t.Conn.Write(packedData)
//	if err != nil {
//		t.Counter = t.Counter - 1
//		return err
//	}
//	return nil
//}