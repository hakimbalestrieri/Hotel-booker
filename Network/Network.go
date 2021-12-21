package Network

/**
File: Network.go
Authors: Hakim Balestrieri
Date: 13.11.2021
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
	"time"
)

type Network struct {
	CurrentServerId int
	OtherServers    map[int]*net.Conn
	Ready           chan int
	//RaymondMsg      chan p.RaymondProtocol
	Raymond      *Raymond
}

// ConnectToServers permet d'ouvrir une connexion sur tous les autres servers du pool
func (network *Network) ConnectToServers() bool {

	network.OtherServers = make(map[int]*net.Conn, conf.ServerNumber-1)
	for i := 0; i < conf.ServerNumber-1; i++ {
		network.OtherServers[i] = nil
	}

	connOK := make(chan bool)
	go network.openConnection(connOK, conf.ServerPorts[network.CurrentServerId])

	for {
		for servId, servPort := range conf.ServerPorts {
			// On ne refait pas une connexion si la connexion a déjà été stockée
			if servId > network.CurrentServerId &&
				network.OtherServers[network.getRelativeIndex(servId)] == nil {
				log.Println(fmt.Sprintf("Tentative de connexion à localhost:%d", servPort))
				conn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", servPort))
				if err == nil {
					conn.Write([]byte("/hello " + strconv.Itoa(network.CurrentServerId) + "\n"))
					log.Println(fmt.Sprintf("Connexion réussi à localhost:%d", servPort))
					network.OtherServers[network.getRelativeIndex(servId)] = &conn
					go network.listen(&conn, servId)
				} else {
					log.Println(fmt.Sprintf("Connexion échouée à localhost:%d", servPort))
				}
			}
		}
		connectedToAll := true
		for servId, v := range network.OtherServers {
			if v == nil && network.getRealIndex(servId) > network.CurrentServerId {
				connectedToAll = false
			}
		}
		if connectedToAll {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	done := <-connOK

	for _, conn := range network.OtherServers {
		_, err := (*conn).Write([]byte("/ready\n"))
		if err != nil {
			return false
		}
	}

	return done
}

// openConnection permet d'écouter si les autres serveurs du pool essaie de se connecter sur ce server
func (network *Network) openConnection(connOk chan bool, port int) {

	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		connOk <- false
		log.Printf("unable to listen on port %d: with connOk : %s", conf.ServerPorts[network.CurrentServerId], err.Error())
		return
	}
	log.Printf("Ouverture du port : %d", conf.ServerPorts[network.CurrentServerId])

	defer listener.Close()

	for {
		connectedToAll := true
		for servId, v := range network.OtherServers {
			if v == nil && network.getRealIndex(servId) < network.CurrentServerId {
				connectedToAll = false
			}
		}
		if connectedToAll {
			break
		}
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept connection: %s", err.Error())
			connOk <- false
			break
		} else {
			log.Printf("Listen to message /hello")

			msg, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				return
			}
			msg = strings.Trim(msg, "\r\n")
			args := strings.Split(msg, " ")
			cmd := strings.TrimSpace(args[0])
			if cmd == "/hello" {
				servId, _ := strconv.Atoi(args[1])
				log.Println(fmt.Sprintf("Server %d : connexion établie", servId))
				network.OtherServers[network.getRelativeIndex(servId)] = &conn
				go network.listen(&conn, servId)
			}
		}
	}
	connOk <- true
	return
}

// listen for any request that are sent to this server
func (network *Network) listen(conn *net.Conn, id int) {
	log.Println("Listenning to server " + strconv.Itoa(id))
	for {
		msg, err := bufio.NewReader(*conn).ReadString('\n')
		if err != nil {
			return
		}
		msg = strings.Trim(msg, "\r\n")
		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		arguments := args[1:]
		if cmd == "/ready" {
			log.Println("Server " + strconv.Itoa(id) + " is ready")
			network.Ready <- id
		} else {
			network.reqLamport(cmd, arguments, id)
			network.reqUpdate(cmd, arguments, id)
		}
	}
}

/*
Store the list of children of each node And we send updates to all nodes
*/
func (network *Network) reqUpdate(cmd string, arguments []string, id int) {
	if cmd == "/update" {
		log.Println("Received update request")
		var update p.Update
		update.NodeId = id
		update.Children = make([]int, 0)
		for _, arg := range arguments {
			child, _ := strconv.Atoi(arg)
			update.Children = append(update.Children, child)
		}
		network.Update <- update
	}
}

// SendUpdateBroadcast goroutine permetant d'envoyer sur les autres serveurs dès qu'une update est disponible dans le channel dédié
func (network *Network) SendUpdateBroadcast() {
	for {
		log.Println("Read msg from UpdateMsgBroadcast")
		message := <-network.UpdateMsgBroadcast
		for i, _ := range network.OtherServers {
			switch message.ReqType {
			case p.UPD_CLIENT:
				_, err := (*network.OtherServers[i]).Write([]byte("/upt_client " + message.Arguments[0] + "\n"))
				if err != nil {
					log.Println("Error during update client")
					continue
				}
			case p.UPD_ROOM:
				_, err := (*network.OtherServers[i]).Write([]byte("/upt_rooms " + message.Arguments[0] + " " + message.Arguments[1] + " " + message.Arguments[2] + " " + message.Arguments[3] + "\n"))
				if err != nil {
					log.Println("Error during update rooms")
					continue
				}
			default:
			}
		}
	}
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

//sendToRootNode(msg, emetteurId) --> 

//send message to root node
func (network *Network) sendToRootNode(msg string, emetteurId int) {
	log.Println("Send message to root node")
	for i, _ := range network.OtherServers {
		if i == 0 {
			_, err := (*network.OtherServers[i]).Write([]byte(msg + "\n"))
			if err != nil {
				log.Println("Error during sendToRootNode")
				continue
			}
		}
	}
}

//send message to every child nodes contained in raymond.queue
func (network *Network) sendToChildNodes(msg string, emetteurId int) {
	log.Println("Send message to child nodes")
	for i, _ := range Raymond.queue {
		if i != 0 {
			_, err := (*network.OtherServers[i]).Write([]byte(msg + "\n"))
			if err != nil {
				log.Println("Error during sendToChildNodes")
				continue
			}
		}
	}
}
				

//sendToChild(msg, emetteurId)

}
