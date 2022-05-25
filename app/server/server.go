package server

import (
	protocol2 "Pangolin/core/protocol"
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"
)

func ServerStart(port string) {
	localAddress,err := net.ResolveTCPAddr("tcp",fmt.Sprintf(":%s",port))
	if err != nil {
		panic(err)
	}
	listener, err := net.ListenTCP("tcp", localAddress)
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for true {
		conn,err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			conn := conn // !

			sessionKey, err := protocol2.HandShakeHander(conn)
			if err != nil {
				return 
			}
			tunnel := protocol2.Tunnel{
				SessionKey: sessionKey,
				Counter:    0,
				Conn:       conn,
				Mutex: sync.Mutex{},
				Buffer: map[uint64][]byte{},
			}

			err = handler(&tunnel)
			if err != nil {
				fmt.Printf("server handle tunnel error: %s\n",err)
			}
		}()
	}
}

func handler(tunnel *protocol2.Tunnel) error {
	// http parse
	n,data,err := tunnel.Read()
	if err != nil {
		fmt.Printf("http parse read error: %s\n", err)
	}
	if n == 0 {
		//return handler(client)
		return errors.New("http parse get null data received")
	}
	//data := protocol.Decrypt(buf[:n], tunnel.SessionKey)
	reader := bufio.NewReader(bytes.NewReader(data))
	request, err := http.ReadRequest(reader)
	if err != nil {
		fmt.Printf("http parse error %s\n", err)
		return err
	}
	host := ""
	if request.URL.Port() == "" {
		host = fmt.Sprintf("%s:80", request.Host)
	} else {
		host = request.Host
	}
	fmt.Printf("proxy to %s :) \n", host)

	// remote
	destConn, err := net.DialTimeout("tcp", host, 3*time.Second)
	if err != nil {
		fmt.Printf("connect remote target server error %s\n", err)
		return err
	}
	if request.URL.Scheme == "http" {
		destConn.Write(data)
	} else {
		err = tunnel.Write([]byte("HTTP/1.1 200 Connection established\r\n\r\n"))
		if err != nil {
			fmt.Printf("send https ok error : %s\n", err)
			return err
		}
	}

	defer destConn.Close()
	defer tunnel.Conn.Close()

	var finished = make(chan bool)
	//var wg = sync.WaitGroup{}
	//wg.Add(1)
	//var mx = sync.Mutex{}
	//mx.Lock()
	// client -> target
	go func() {
		//defer wg.Done()
		//defer mx.Unlock()
		for true {
			n, data, err := tunnel.Read()
			if err != nil {
				fmt.Printf("client -> target tunnel read error: %s\n", err)
				break
			}

			if n != 0 {
				_, err = destConn.Write(data)
				if err != nil {
					fmt.Printf("client -> target tunnel write error: %s\n", err)
					break
				}
			}
		}
		finished <- true
	}()

	// target -> client
	go func() {
		//defer wg.Done()
		//defer mx.Unlock()
		var buf = make([]byte,1<<16)
		for true {
			n, err := destConn.Read(buf)
			if err != nil {
				fmt.Printf("target -> client conn read error: %s\n", err)
				break
			}
			if n != 0 {
				err = tunnel.Write(buf[:n])
				if err != nil {
					fmt.Printf("target -> client conn write error: %s\n", err)
					break
				}
			}
		}
		finished <- true
	}()

	<-finished
	return nil
}
