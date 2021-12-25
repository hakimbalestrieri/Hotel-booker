package Mutex

/**
File: Mutex.go
Authors: Hakim Balestrieri
Date: 25.12.2021
*/

import (
	"PRR-Labo3-Balestrieri/Business"
	"PRR-Labo3-Balestrieri/Config"
	"PRR-Labo3-Balestrieri/Network"
	p "PRR-Labo3-Balestrieri/Protocol"
	Raymond "PRR-Labo3-Balestrieri/Raymond"
	"log"
	"net"
)

type Mutex struct {
	ServerId           int
	network            *Network.Network
	hotel              *Business.Hotel
	Raymond            *Raymond.Raymond
	MsgSC              chan string
	AccessCS           chan bool
	UserMsgIn          chan Network.UserProtocol
	RaymondMsgIn       chan p.RaymondProtocol
	RaymondMsgOut      chan p.RaymondProtocol
	UpdateMsgIn        chan p.UpdateProtocol
	UpdateMsgOut chan p.UpdateProtocol
}

func NewMutex(id int) *Mutex {
	return &Mutex{ServerId: id}
}

func (mutex *Mutex) Init() {

	mutex.RaymondMsgIn = make(chan p.RaymondProtocol)
	mutex.RaymondMsgOut = make(chan p.RaymondProtocol)
	mutex.UpdateMsgIn = make(chan p.UpdateProtocol)
	mutex.UpdateMsgOut = make(chan p.UpdateProtocol)
	mutex.UserMsgIn = make(chan Network.UserProtocol)
	mutex.AccessCS = make(chan bool)

	mutex.network = &Network.Network{
		CurrentServerId: mutex.ServerId,
		UpdateMsgIn:     mutex.UpdateMsgIn,
		UpdateMsgOut:    mutex.UpdateMsgOut,
		RaymondMsgIn:    mutex.RaymondMsgIn,
		RaymondMsgOut:   mutex.RaymondMsgOut,
		HasAccess:       make(chan bool),
	}

	if mutex.network.ConnectToRootServer() != true {
		return
	}

	serversSchema := Config.ServerSchema[mutex.ServerId]

	mutex.Raymond = &Raymond.Raymond{
		//canaux
		RayMsgIn:  mutex.RaymondMsgIn,
		RayMsgOut: mutex.RaymondMsgOut,

		//variables raymond
		CurrentId: mutex.ServerId,
		Status:    p.RAY_NO,
		ParentId:  serversSchema.Root,
	}

	mutex.MsgSC = make(chan string)

	// Création de l'Hotel
	rooms := make(map[int]map[int]string)
	mutex.hotel = &Business.Hotel{
		Network:   mutex.network,
		Rooms:     rooms,
		MsgSC:     mutex.MsgSC,
		AccessCS:  mutex.AccessCS,
		UserMsgIn: mutex.UserMsgIn,
	}

	for i := 0; i < Config.RoomsNumber; i++ {
		mutex.hotel.Rooms[i] = make(map[int]string)
	}
}

// Run launch goroutines
func (mutex *Mutex) Run() {
	go mutex.hotel.Run()
	go mutex.listenReqSC()
	go mutex.network.ManageOutMsg()
	go mutex.network.ManageUpdateMsgOut()
	go mutex.updateHotel()
	go mutex.Raymond.Run()
}

// NewClient create a new Client and listen to tcp Client
// used as a goroutine
func (mutex *Mutex) NewClient(conn net.Conn) {

	log.Printf("new Client has joined: %s", conn.RemoteAddr().String())

	c := &Network.Client{
		Conn:     conn,
		Username: "newClient",
		Commands: mutex.UserMsgIn,
	}

	defer c.Conn.Close()

	c.ReadInput()
}

// updateHotel goroutine that takes care of updating the hotel everytime we recieve a request on the dedicated channel
func (mutex *Mutex) updateHotel() {
	for {
		log.Println("Read msg from UpdateMsgIn")
		message := <-mutex.UpdateMsgIn
		switch message.ReqType {
		case p.UPD_CLIENT:
			mutex.hotel.UpdateUsername(message.Arguments[0])
		case p.UPD_ROOM:
			mutex.hotel.UpdateRooms(message.Arguments)
		default:
		}
	}
}

// listenReqSC listen channel for incoming critical section request
func (mutex *Mutex) listenReqSC() {
	for {
		log.Println("Read msg from MsgSC")
		msg := <-mutex.MsgSC
		switch msg {
		case "req":
			log.Println("Put msg REQ in RayMsgIn")
			mutex.RaymondMsgIn <- p.RaymondProtocol{
				ReqType: p.RAYMOND_PRO_REQ,
			}
		case "rel":
			log.Println("Put msg REL in RayMsgIn")
			mutex.RaymondMsgIn <- p.RaymondProtocol{
				ReqType: p.RAYMOND_END,
			}
		default:
		}
	}
}
