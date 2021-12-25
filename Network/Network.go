package Network

/**
File: Network.go
Authors: Hakim Balestrieri
Date: 25.12.2021
*/

import (
	conf "PRR-Labo3-Balestrieri/Config"
	p "PRR-Labo3-Balestrieri/Protocol"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
)

type Network struct {
	CurrentServerId    int
	UpdateMsgIn        chan p.UpdateProtocol
	UpdateMsgBroadcast chan p.UpdateProtocol
	rootNode           *Server
	childNodes         []*Server
}

// ConnectToRootServer permet d'ouvrir une connexion sur le serveur root
func (network *Network) ConnectToRootServer() bool {
	if conf.Debug {
		log.Println(fmt.Sprintf("Launching ConnectToRootServer() ..."))
	}
	serversSchema := conf.ServerSchema[network.CurrentServerId]
	network.childNodes = make([]*Server, len(serversSchema.Children), len(serversSchema.Children))
	initConnection := make(chan bool)

	go network.ConnectToChildren(initConnection)

	if serversSchema.Root != network.CurrentServerId {
		rootNode := conf.ServerPorts[serversSchema.Root]
		if conf.Debug {
			log.Println(fmt.Sprintf("Connecting to : %d", rootNode))
		}
		for {
			conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", rootNode))
			if err != nil {
				log.Println(fmt.Sprintf("Failed to connect to node : %d", rootNode))
			} else {
				conn.Write([]byte("/hello " + strconv.Itoa(network.CurrentServerId) + "\n"))
				log.Println(fmt.Sprintf("Sucessfully connected to node : %d", rootNode))
				network.rootNode = &Server{
					id:   serversSchema.Root,
					conn: &conn,
				}
				go network.listen(network.rootNode)
				break
			}
		}

	}

	return <-initConnection
}

func (network *Network) ConnectToChildren(connOk chan bool) {

	if conf.Debug {
		log.Println("Launching ConnectToChildren")
	}

	serverPort := conf.ServerPorts[network.CurrentServerId]

	listener, err := net.Listen("tcp", ":" + strconv.Itoa(serverPort))
	if err != nil {
		connOk <- false
		log.Printf("unable to listen on port %d: with connOk : %s\n", conf.ServerPorts[network.CurrentServerId], err.Error())
		return
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	if conf.Debug {
		log.Printf("Connected to node : %d\n", serverPort)
	}

	for {
		connectedToAll := true
		for _, childConnection := range network.childNodes {
			if childConnection == nil {
				connectedToAll = false
			}
		}
		if connectedToAll {
			break
		}

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s\n", err.Error())
			connOk <- false
			break
		} else {
			if conf.Debug {
				log.Println("Listen to message /hello")
			}

			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				if conf.Debug {
					log.Printf("Error while reading string : %s\n", err.Error())
				}
				connOk <- false
				return
			}

			msg = strings.Trim(msg, "\r\n")
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])
			if cmd == "/hello" {
				servId, _ := strconv.Atoi(args[1])
				if conf.Debug {
					log.Printf("Connected to node : %d\n", servId)
				}
				network.childNodes = append(network.childNodes, &Server{
					id:   servId,
					conn: &conn,
				})
			}
		}
	}

	if conf.Debug {
		log.Println("ConnectToChildren terminated")
	}
	connOk <- true
	return
}

// listen for any request that are sent to this server
func (network *Network) listen(server *Server) {
	log.Println("Listening to Server " + strconv.Itoa(server.id))
	for {
		msg, err := bufio.NewReader(*server.conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		arguments := args[1:]
		network.reqUpdate(cmd, arguments, server.id)
	}
	log.Println("End listening to Server " + strconv.Itoa(server.id))
}

// reqUpdate permet d'envoyer sur le channel dédié une demande d'update
func (network *Network) reqUpdate(cmd string, arguments []string, relId int) {
	switch cmd {
	case "/upt_client":
		log.Println("Put msg /upt_client in UpdateMsgIn")
		network.UpdateMsgIn <- p.UpdateProtocol{
			ReqType:      p.UPD_CLIENT,
			Arguments:    arguments,
			ServerIdFrom: relId,
		}
	case "/upt_rooms":
		log.Println("Put msg /upt_rooms in UpdateMsgIn")
		network.UpdateMsgIn <- p.UpdateProtocol{
			ReqType:      p.UPD_ROOM,
			Arguments:    arguments,
			ServerIdFrom: relId,
		}
	default:
		log.Println("Could not recognized reqUpdate")
	}
}
