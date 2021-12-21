package Mutex2


import (
	p "PRR-Labo3-Balestrieri/Protocol"
	"PRR-Labo3-Balestrieri/Network"
)

type Mutex struct {
	network *Network.Network
	hotel	*Business.Hotel
	raymond *Raymond.Raymond
	UserMsgIn	chan Network.UserProtocol
	UpdateMsgIn	chan p.UpdateProtocol
	UpdateMsgBroadcast	chan p.UpdateProtocol
}

func NewMutex(id int) *Mutex {
	mutex := &Mutex{
		network: Network.NewNetwork(id),
		hotel: Business.NewHotel(id),
		raymond: Raymond.NewRaymond(id),
		UserMsgIn: make(chan Network.UserProtocol),
		UpdateMsgIn: make(chan p.UpdateProtocol),
		UpdateMsgBroadcast: make(chan p.UpdateProtocol),
	}
	
	connected := mutex.network.ConnectToServers()

	if connected != true {
		return
	}

	//initialize mutex.raymond
	mutex.raymond = &Raymond.Raymond{
		me: id,
		status: p.RAY_NO,
		parent: nil,
		file: 0,
		RaymondMsgBroadcast: make(chan p.RaymondProtocol),
	}

	mutex.waitForEveryone()
	mutex.raymond.Init()

	// Create hotel
	rooms := make(map[int]map[int]string)
	mutex.hotel = &Business.Hotel{
		Network: 	mutex.network,
		Rooms:   	rooms,
		MsgSC:   	mutex.MsgSC,
		AccessCS: 	mutex.AcessCS,
		UserMsgIn: 	mutex.UserMsgIn,
	}

	for i := 0; i < Config.RoomsNumber; i++ {
		mutex.hotel.Rooms[i] = make(map[int]string)
	}

}

//Run goroutines
func (mutex *Mutex) Run() {
	go mutex.raymond.Run()
	go mutex.hotel.Run()
	go mutex.updateHotel()
}

//update hotel username or update room everytime we recevie update message
func (mutex *Mutex) updateHotel() {
	for {
		select {
		case msg := <-mutex.UpdateMsgIn:
			switch msg.Type {
			case p.UPDATE_USER:
				mutex.hotel.UpdateUsername(msg.Arguments[0])
			case p.UPDATE_ROOM:
				mutex.hotel.UpdateRooms(msg.Arguments)
			}
		}
	}


/*
 create new client and listen to tcp client
*/
func (ray *Raymond) NewClient(conn net.Conn) {
	log.Printf("new Client has joined: %s", conn.RemoteAddr().String())

	ray.RaymondMsgBroadcast <- RaymondProtocol{
		ReqType:  RAY_REQ,
		ServerId: ray.me,
		ParentId: ray.parent,
	}
	//defer close connection
	defer conn.Close()

	//read data from client
	for {
		data := make([]byte, 1024)
		length, err := conn.Read(data)
		if err != nil {
			log.Printf("error: %v", err)
			break
		}
		//log.Printf("read data: %s", data[:length])
		ray.RaymondMsgBroadcast <- RaymondProtocol{
			ReqType:  RAY_REQ,
			ServerId: ray.me,
			ParentId: ray.parent,
		}
	}

}

/*
wait for all servers to connect
*/
func (mutex *Mutex) waitForEveryone() {
	var array []bool
	array = make([]bool, Config.ServerNumber)
	array[mutex.ServerId] = true
	for {
		ok := true
		for _, val := range array {
			if !val {
				ok = false
				break
			}
		}
		if ok {
			break
		}
		msg := <- mutex.network.Ready
		array[msg] = true
	}
	log.Println("All servers are connected")
}