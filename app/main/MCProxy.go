package main

import (
	"Pangolin/app/client"
	"Pangolin/app/server"
	"Pangolin/core/config"
	"fmt"
	"gopkg.in/alecthomas/kingpin.v2"
	"io/ioutil"
	"os"
)

var (
	app = kingpin.New("MCProxy", "A Simple TCP Based Http Proxy.")
	// vpn client 1234 192.168.6.131:1234
	clientMode = app.Command("client", "client mode")
	localPort = clientMode.Arg("port", "client local port").Required().String()
	serverAddress = clientMode.Arg("server", "server address:port. (like 192.168.6.131:1234)").Default("1234").String()
	clientConfig = clientMode.Arg("configPath","config path, default ./client.json").Default("./client.json").String()
	// server mode
	serverMode = app.Command("server", "server mode")
	serverPort = serverMode.Arg("port", "port to accept client's connection").Default("4321").String()
	serverConfig = serverMode.Arg("configPath","config path, default ./server.json").Default("./server.json").String()
)

func main() {
	var isServer = false
	var configPath = ""
	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case clientMode.FullCommand():
		isServer = false
		fmt.Println("[+] client port: " + *localPort)
		fmt.Println("[+] server address: " + *serverAddress)
		fmt.Println("config: " + *clientConfig)
		configPath = *clientConfig
	case serverMode.FullCommand():
		isServer = true
		fmt.Println("[+] server port: " + *serverPort)
		fmt.Println("config: " + *serverConfig)
		configPath = *serverConfig
	}

	jsonFile, err := os.Open(configPath)
	defer jsonFile.Close()
	if err != nil {
		fmt.Printf("open config file error %s",err)
		return
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	fmt.Println(string(byteValue))

	err = config.InitConfig(byteValue)
	if err != nil {
		fmt.Printf("config load error %s",err)
		return
	}

	if isServer {
		server.ServerStart(*serverPort)
	} else {
		client.ClientStart(*serverAddress, *localPort)
	}
}