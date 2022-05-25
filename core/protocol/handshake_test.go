package protocol

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestHandshake(t *testing.T) {
	go func() {
		// server
		addr,err := net.ResolveTCPAddr("tcp",":1234")
		if err != nil {
			panic(err)
		}
		listener,err := net.ListenTCP("tcp",addr)
		//listener.SyscallConn()
		if err != nil {
			panic(err)
		}
		defer listener.Close()
		fmt.Println("server listening in 0.0.0.0:1234")
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}
		sessionKey1, err := HandShakeHander(conn)
		if err != nil {
			panic(err)
		}
		fmt.Printf("server session key: %s\n", sessionKey1)

		var counter byte = 0
		for true {
			conn := conn
			var buf = make([]byte,1024)
			n,err := conn.Read(buf)
			if err != nil {
				panic(err)
			}
			if n != 0 {
				counter = counter + 1
				data, err := UnPackData(buf[:n], sessionKey1,counter)
				if err != nil {
					panic(err)
				}
				counter = counter + 1
				fmt.Printf("%s\n", data)
				ret := []byte("hello client!")
				packedData, err := PackData(ret,sessionKey1,counter)
				if err != nil {
					panic(err)
				}
				_, err = conn.Write(packedData)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	tunnel, err := GetConnect("127.0.0.1:1234")
	if err != nil {
		panic(err)
	}
	//defer conn.Close()
	fmt.Printf("client session key: %s\n", tunnel.SessionKey)
	time.Sleep(1*time.Second)

	var clientCounter byte = 0
	clientCounter = clientCounter + 1
	packedData, err := PackData([]byte("hello server"),tunnel.SessionKey,clientCounter)
	if err != nil {
		panic(err)
	}
	_, err = tunnel.Conn.Write(packedData)
	if err != nil {
		panic(err)
	}

	for true {
		conn := tunnel.Conn
		var buf = make([]byte,1024)
		n,err := conn.Read(buf)
		if err != nil {
			panic(err)
		}
		if n != 0 {
			clientCounter = clientCounter + 1
			data, err := UnPackData(buf[:n], tunnel.SessionKey,clientCounter)
			if err != nil {
				panic(err)
			}
			clientCounter = clientCounter + 1
			fmt.Printf("%s\n", data)
			ret := []byte("hello server!")
			packedData, err := PackData(ret,tunnel.SessionKey,clientCounter)
			if err != nil {
				panic(err)
			}
			_, err = conn.Write(packedData)
			if err != nil {
				panic(err)
			}
		}
	}
}
