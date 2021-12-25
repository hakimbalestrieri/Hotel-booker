/**
File: tcpServer.go
Authors: Hakim Balestrieri & Alexandre Mottier
Date: 22.10.2021
*/
package main

import (
	"PRR-Labo3-Balestrieri/Config"
	"PRR-Labo3-Balestrieri/Mutex"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {

	servID := os.Args[1]
	serverId, err := strconv.Atoi(strings.Trim(servID, "\r\n"))

	if err != nil || serverId < 0 || serverId > Config.ServerNumber - 1 {
		log.Fatalf("Build the tcpServer with his number (from 1 to %d): %s", Config.ServerNumber, err.Error())
	}

	debug := false

	if len(os.Args) > 2 {
		if os.Args[2] == "1" {
			debug = true
		}
	}

	Config.Debug = debug

	mutex := Mutex.NewMutex(serverId)

	mutex.Init()

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(Config.ClientPorts[mutex.ServerId]))
	if err != nil {
		log.Fatalf("unable to start hotel: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started")
	mutex.Run()
	log.Printf("Starting goroutines")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			continue
		}

		go mutex.NewClient(conn)
	}
}
