/**
File: tcpClient.go
Authors: Hakim Balestrieri
Date: 22.10.2021
*/
package main

import (
	"PRR-Labo3-Balestrieri/Config"
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	arguments := os.Args
	if len(arguments) == 1 {
		fmt.Println("Please provide server ID.")
		return
	}

	CONNECT, _ := strconv.Atoi(arguments[1])
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(Config.ClientPorts[CONNECT-1]))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	clientReader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)

	for {
		// Waiting for the Client request
		clientRequest, err := clientReader.ReadString('\n')

		switch err {
		case nil:
			clientRequest := strings.TrimSpace(clientRequest)
			if _, err = conn.Write([]byte(clientRequest + "\n")); err != nil {
				log.Printf("failed to send the Client request: %v\n", err)
			}
		case io.EOF:
			log.Println("Client closed the connection")
			return
		default:
			log.Printf("Client error: %v\n", err)
			return
		}

		// Waiting for the ServerId response
		serverResponse, err := serverReader.ReadString('\n')

		switch err {
		case nil:
			log.Println(strings.TrimSpace(serverResponse))
		case io.EOF:
			log.Println("ServerId closed the connection")
			return
		default:
			log.Printf("ServerId error: %v\n", err)
			return
		}
	}
}
