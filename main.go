package main

import (
	"net"
	"os"
	"github.com/vimukthi-git/beanstalkg/operation"
	"encoding/json"
	"log"
	"github.com/vimukthi-git/beanstalkg/architecture"
)

func main() {
	service := ":11300"
	tubeRegister := make(chan architecture.Command)
	// use this tube to send the channels for each individual tube to the clients when the do 'use' command
	tubeHandlers := make(chan chan architecture.Command)
	stop := make(chan bool)
	operation.NewTubeRegister(tubeRegister, tubeHandlers, stop)
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	for {
		log.Println("Waiting..")
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		operation.NewClientHandler(conn, tubeRegister, tubeHandlers, stop)
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal("Fatal error:", err.Error())
	}
}

type Configuration struct {
	Beanstalks []string `json:"beanstalks"`
}

func getConfig(env string) Configuration {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	configuration := make(map[string]Configuration)
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("error in parsing config:", err)
	}
	envConf, ok := configuration[env]
	if !ok {
		log.Fatal("No configuration found for the given environment name")
	}
	return envConf
}