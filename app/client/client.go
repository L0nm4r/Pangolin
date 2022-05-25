package client

import (
	"Pangolin/core/protocol"
	"fmt"
	"io"
	"net"
)

func ClientStart(remoteAddress string, localPort string) {
	// local client listening address
	addr,err := net.ResolveTCPAddr("tcp",fmt.Sprintf(":%s", localPort))
	if err != nil {
		panic(err)
	}

	// bind to listening
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}

	defer listener.Close()

	// waiting for connections
	for true {
		localClient,err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go func() {
			localClient := localClient // !!
			defer localClient.Close()

			tunnel, err := protocol.GetConnect(remoteAddress)
			defer tunnel.Conn.Close()
			if err != nil {
				fmt.Printf("get connect error: %s\n",err)
				return
			}

			var finished = make(chan bool)
			//var wg = sync.WaitGroup{}
			//wg.Add(2)
			//var mx = sync.Mutex{}
			//mx.Lock()
			// remote server -> client
			go func() {
				//defer wg.Done()
				//defer mx.Unlock()
				for true {
					n, data, err := tunnel.Read()
					if err != nil {
						if err != io.EOF {
							fmt.Printf("server -> client'tunnel read error: %s\n",err)
						}
						break
					}
					if n != 0 {
						_, err = localClient.Write(data)
						if err != nil {
							if err != io.EOF {
								fmt.Printf("server -> client'tunnel write error: %s\n",err)
							}
							break
						}
					}
				}
				finished <- true
			}()

			// client -> remote server
			go func() {
				//defer wg.Done()
				//defer mx.Unlock()
				var buf = make([]byte,1<<16)
				for true {
					n, err := localClient.Read(buf)
					if err != nil {
						if err != io.EOF {
							fmt.Printf("client -> remote'tunnel read error: %s\n",err)
						}
						break
					}
					if n != 0 {
						err = tunnel.Write(buf[:n])
						if err != nil {
							if err != io.EOF {
								fmt.Printf("client -> remote'tunnel write error: %s\n",err)
							}
							break
						}
					}
				}
				finished <- true
			}()
			//wg.Wait()
			//mx.Lock()
			//mx.Unlock()
			<- finished
		}()
	}
}
