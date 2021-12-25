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
	//Canaux
	UpdateMsgIn   chan p.UpdateProtocol
	UpdateMsgOut  chan p.UpdateProtocol
	RaymondMsgIn  chan p.RaymondProtocol
	RaymondMsgOut chan p.RaymondProtocol
	HasAccess     chan bool

	//Variables Network
	CurrentServerId int
	rootNode        *Server
	childNodes      []*Server
}

// ConnectToRootServer permet d'ouvrir une connexion sur le serveur root
func (network *Network) ConnectToRootServer() bool {
	if conf.Debug {
		log.Println(fmt.Sprintf("Launching ConnectToRootServer() ..."))
	}
	serversSchema := conf.ServerSchema[network.CurrentServerId]
	network.childNodes = make([]*Server, len(serversSchema.Children), len(serversSchema.Children))
	initConnection := make(chan bool)

	go network.ConnectToChildServers(initConnection)

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
				network.WriteToNodeWithParameter(network.CurrentServerId, "/hello", strconv.Itoa(network.CurrentServerId))
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

func (network *Network) ConnectToChildServers(connOk chan bool) {

	if conf.Debug {
		log.Println("Launching ConnectToChildServers")
	}

	serverPort := conf.ServerPorts[network.CurrentServerId]

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(serverPort))
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
				childServer := &Server{
					id:   servId,
					conn: &conn,
				}

				network.childNodes = append(network.childNodes, childServer)
				go network.listen(childServer)
			}
		}
	}

	if conf.Debug {
		log.Println("ConnectToChildServers terminated")
	}
	connOk <- true
	return
}

func (network *Network) listen(serv *Server) {

	log.Printf("Listen server launched, server : %d\n", serv.id)

	var delimiter byte = '\n'
	cutset := "\r\n"
	separator := " "

	for {
		msg, err := bufio.NewReader(*serv.conn).ReadString(delimiter)
		if err != nil {
			return
		}
		msg = strings.Trim(msg, cutset)
		args := strings.Split(msg, separator)
		cmd := strings.TrimSpace(args[0])
		argument := args[1:]

		switch cmd {

		case "/token":
			log.Printf("Server handle command : %s", cmd)
			network.RaymondMsgIn <- p.RaymondProtocol{
				ReqType: p.RAYMOND_REQ,
			}
		case "/req":
			serverId, _ := strconv.Atoi(argument[0])
			log.Printf("Server with id : %d handle command : %s", serverId, cmd)
			network.RaymondMsgIn <- p.RaymondProtocol{
				ReqType:  p.RAYMOND_REQ,
				ServerId: serverId,
			}
		case "/ready":
			log.Printf("Server with id : %d handle command : %s", serv.id, cmd)
			serv.hasAccess = true
		case "/go":
			log.Printf("Server with id : %d handle command : %s", serv.id, cmd)
			network.HasAccess <- true
		case "/upt_client":
			log.Printf("Server handle command : %s", cmd)
			network.UpdateMsgIn <- p.UpdateProtocol{
				ReqType:      p.UPD_CLIENT,
				Arguments:    argument,
				ServerIdFrom: serv.id,
			}
		case "/upt_rooms":
			log.Printf("Server handle command : %s", cmd)
			network.UpdateMsgIn <- p.UpdateProtocol{
				ReqType:      p.UPD_ROOM,
				Arguments:    argument,
				ServerIdFrom: serv.id,
			}
		default:
			if conf.Debug {
				log.Println("Default case listen()")
			}
		}
	}
}

func (network *Network) waitUntilServersAreReady() bool {

	if conf.Debug {
		log.Println("Waiting until servers are ready")
	}

	serversSchema := conf.ServerSchema[network.CurrentServerId]
	childLength := len(serversSchema.Children)
	if childLength == 0 {
		//Le serveur n'a aucun enfant , dans ce cas-là il passe directement à ready
		network.WriteToNode(network.rootNode.id, "/ready")
	} else {
		childNodes := make([]bool, childLength)
		for {
			for i, server := range network.childNodes {
				if server.hasAccess == true {
					childNodes[i] = true
				}
			}
			everyNodesAreReady := true

			for _, elem := range childNodes {
				if elem == false {
					everyNodesAreReady = false
				}
			}
			if everyNodesAreReady {
				break
			}
		}
		if conf.Debug {
			log.Println("All nodes are ready")
		}
	}

	rootNodeId := network.rootNode.id

	if rootNodeId != network.CurrentServerId {
		network.WriteToNode(network.rootNode.id, "/ready")
		<-network.HasAccess
		log.Printf("Server with id %d is ready \n", rootNodeId)

	}

	for _, server := range network.childNodes {
		network.WriteToNode(server.id, "/go")
	}
	return true
}

func (network *Network) ManageOutMsg() {
	for {
		if conf.Debug {
			log.Println("Currently reading msg from RaymondMsgOut")
		}

		msg := <-network.RaymondMsgOut
		nodeId := network.rootNode.id

		switch msg.ReqType {

		case p.RAYMOND_TOKEN:

			var servers []*Server
			if network.rootNode.id == network.CurrentServerId {
				for _, serv := range network.childNodes {
					if serv.id != msg.ServerId {
						servers = append(servers, serv)
					} else {
						network.rootNode = serv
					}
				}
				network.childNodes = servers
			}
			network.WriteToNode(nodeId, "/token")

		case p.RAYMOND_REQ:
			network.WriteToNodeWithParameter(nodeId, "/req", strconv.Itoa(network.CurrentServerId))

		default:
			if conf.Debug {
				log.Println("Default case, ManageOutMsg")
			}
		}
	}
}

func (network *Network) WriteToNode(nodeId int, typeMsg string) {
	_, err := (*network.rootNode.conn).Write([]byte(typeMsg + "\n"))
	if err != nil {
		log.Println(fmt.Sprintf("Server with id  %d failed to send "+typeMsg+" message", nodeId))
	} else {
		log.Println(fmt.Sprintf("Message sent to server with id : %d", nodeId))
	}
}

func (network *Network) WriteToNodeWithParameter(nodeId int, typeMsg string, query string) {
	_, err := (*network.rootNode.conn).Write([]byte(typeMsg + " " + query + "\n"))
	if err != nil {
		log.Println(fmt.Sprintf("Server with id :  %d failed to send message : %s", nodeId, typeMsg))
	} else {
		log.Println(fmt.Sprintf("Message sent to server with id : %d", nodeId))
	}
}

func (network *Network) ManageUpdateMsgOut() {
	for {
		if conf.Debug {
			log.Println("Currently reading msg from UpdateMsgOut")
		}
		msg := <-network.UpdateMsgOut
		for _, serv := range network.childNodes {
			switch msg.ReqType {
			case p.UPD_ROOM:

				var arguments []string
				for i := 0; i <= 3; i++ {
					arguments = append(arguments, msg.Arguments[i]+" ")
				}
				arguments = append(arguments, "\n")
				network.WriteToNodeWithParameter(serv.id, "/upt_rooms", strings.Join(arguments, " "))

			case p.UPD_CLIENT:
				network.WriteToNodeWithParameter(serv.id, "/upt_client", msg.Arguments[0])

			default:
				if conf.Debug {
					log.Println("Default case ManageUpdateMsgOut()")
				}
			}
		}
	}
}
