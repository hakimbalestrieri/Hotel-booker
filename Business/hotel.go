/**
File: hotel.go
Authors: Hakim Balestrieri
Date: 22.10.2021
*/
package Business

import (
	"PRR-Labo3-Balestrieri/Config"
	"PRR-Labo3-Balestrieri/Network"
	"PRR-Labo3-Balestrieri/Protocol"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Hotel represent an Hotel with its rooms and commands that can be used
type Hotel struct {
	// clients All clients that register their username
	clients []string
	// Network pointer to the network to send updates
	Network *Network.Network
	// Rooms is a map with first index representing the room number,
	// second index represent the date
	// and contain a string that represent Client username
	Rooms map[int]map[int]string
	// MsgSC Channel to ask for Critical Section
	MsgSC chan string
	// AccessCS Channel to check if hotel can enter in Critical Section
	AccessCS chan bool
	// UserMsgIn channel to communicate between client and hotel
	UserMsgIn chan Network.UserProtocol

	Raymond *Raymond
}

// Run redirect command to the right function to do
func (h *Hotel) Run() {
	log.Printf("Hotel is now running")
	for {
		cmd := <-h.UserMsgIn
		switch cmd.Id {
		case Network.CMD_USERNAME:
			h.Username(cmd.Client, cmd.Args)
		case Network.CMD_RESERVE:
			h.reserve(cmd.Client, cmd.Args)
		case Network.CMD_ROOMS:
			h.listRooms(cmd.Client, cmd.Args)
		case Network.CMD_QUIT:
			h.quit(cmd.Client)
		}
	}
}

// Username affect Username to the Client
func (h *Hotel) Username(c *Network.Client, args []string) {
	//
	if len(args) < 2 {
		log.Printf("Client (%s) tried to choose a username without putting the username", c.Conn.RemoteAddr().String())
		c.Msg("username is required. usage: /username NAME")
		return
	}

	if Config.Debug {
		log.Printf("Client (%s) in username function ask for Critical Section", c.Conn.RemoteAddr().String())
		//time.Sleep(5 * time.Second)
	}
	// Appelle blocant
	h.Request(c)

	if Config.Debug {
		log.Printf("Client (%s) in username function is in Critical Section", c.Conn.RemoteAddr().String())
		//time.Sleep(5 * time.Second)
	}
	// SECTION CRITIQUE
	for _, client := range h.clients {
		if client == args[1] {
			log.Printf("Client (%s) tried to choose a username that already exist", c.Conn.RemoteAddr().String())
			c.Msg("username already exist. usage: /username NAME")
			return
		}
	}
	h.clients = append(h.clients, args[1])
	c.Username = args[1]
	h.Network.UpdateMsgBroadcast <- Protocol.UpdateProtocol{
		ReqType:    Protocol.UPD_CLIENT,
		Arguments:  []string{args[1]},
		ServerIdTo: h.Network.GetIdOtherServers(),
	}
	// FIN SECTION CRITIQUE
	if Config.Debug {
		log.Printf("Client (%s) in username function finish with Critical Section", c.Username)
		//time.Sleep(5 * time.Second)
	}
	h.Release()
	if Config.Debug {
		log.Printf("Client (%s) in username function is out of Critical Section", c.Username)
		//time.Sleep(5 * time.Second)
	}
	log.Printf("Client (%s) choose a username : %s", c.Conn.RemoteAddr().String(), c.Username)
	c.Msg(fmt.Sprintf("Welcome in the hotel, %s", c.Username))
}

// reserve affect a room to the Client if the room is free
func (h *Hotel) reserve(c *Network.Client, args []string) {
	if c.Username == "newClient" {
		log.Printf("Client (%s) tried to reserve without username", c.Conn.RemoteAddr().String())
		c.Msg("You need to be identified, use /username")
		return
	}
	if len(args) < 4 {
		log.Printf("Client (%s) tried to reserve without enough argument", c.Username)
		c.Msg("room name, date of arrival and number of nights are required. usage: /reserve ROOM_NUMBER DATE NUMBER_NIGHTS")
		return
	}
	roomNameRaw := args[1]
	dateReservationRaw := args[2]
	numberNightsRaw := args[3]

	roomName, err := strconv.Atoi(roomNameRaw)
	if err != nil {
		return
	}

	dateReservation, err := strconv.Atoi(dateReservationRaw)
	if err != nil {
		return
	}

	numberNights, err := strconv.Atoi(numberNightsRaw)
	if err != nil {
		return
	}

	// Check for room number is in range
	if roomName < 1 || roomName > Config.RoomsNumber {
		log.Printf("Client (%s) tried to reserve but room number was not in range", c.Username)
		c.Msg("room number was not in range : min 1, max " + strconv.Itoa(Config.RoomsNumber))
		return
	}

	// Check for room number is in range
	if dateReservation < 1 || dateReservation+numberNights > Config.DayNumber {
		log.Printf("Client (%s) tried to reserve but duration was not in range", c.Username)
		c.Msg("duration was not in range : min 1, max " + strconv.Itoa(Config.DayNumber))
		return
	}

	if Config.Debug {
		log.Printf("Client (%s) in reserve function ask for Critical Section", c.Conn.RemoteAddr().String())
		//time.Sleep(5 * time.Second)
	}
	h.Request(c)
	if Config.Debug {
		log.Printf("Client (%s) in reserve function is in Critical Section", c.Conn.RemoteAddr().String())
		//time.Sleep(5 * time.Second)
	}
	// SECTION CRITIQUE
	nok := false
	for date := dateReservation; date < dateReservation+numberNights; date++ {
		if h.Rooms[roomName-1][date-1] != "" {
			nok = true
		}
	}
	// The room is used during the lap of time he wants
	if nok {
		log.Printf("Client (%s) tried to reserve room number: %d but was reserved already", c.Username, roomName)
		c.Msg("The room is already occupied.")
		return
	}

	// assign Client for the room during the number of nights he wants, starting the date he wants
	for date := dateReservation; date < dateReservation+numberNights; date++ {
		h.Rooms[roomName-1][date-1] = c.Username
	}

	// Update other servers
	h.Network.UpdateMsgBroadcast <- Protocol.UpdateProtocol{
		ReqType:    Protocol.UPD_ROOM,
		Arguments:  []string{roomNameRaw, dateReservationRaw, numberNightsRaw, c.Username},
		ServerIdTo: h.Network.GetIdOtherServers(),
	}

	// FIN SECTION CRITIQUE
	if Config.Debug {
		log.Printf("Client (%s) in reserve function finish with Critical Section", c.Username)
		//time.Sleep(5 * time.Second)
	}
	h.Release()
	if Config.Debug {
		log.Printf("Client (%s) in reserve function is out of Critical Section", c.Username)
		//time.Sleep(5 * time.Second)
	}
	log.Printf("Client (%s) reserve room number : %d", c.Username, roomName)
	c.Msg(fmt.Sprintf("Room number: %d is now reserved", roomName))
}

// listRooms show to user the list of Rooms and their availability or show a room that can be used
func (h *Hotel) listRooms(c *Network.Client, args []string) {
	if c.Username == "newClient" {
		log.Printf("Client (%s) tried to list rooms without username", c.Conn.RemoteAddr().String())
		c.Msg("You need to be identified, use /username")
		return
	}
	if len(args) < 2 {
		log.Printf("Client (%s) tried to list rooms without enough argument", c.Username)
		c.Msg("date of arrival is required, number of nights is optional. usage: /rooms DATE (NUMBER_NIGHTS)")
		return
	}

	if len(args) > 3 {
		log.Printf("Client (%s) tried to list rooms with too much arguments", c.Username)
		c.Msg("date of arrival is required, number of nights is optional. usage: /rooms DATE (NUMBER_NIGHTS)")
		return
	}

	dateReservationRaw := args[1]
	dateReservation, err := strconv.Atoi(dateReservationRaw)

	if dateReservation < 1 || dateReservation > Config.DayNumber {
		log.Printf("Client (%s) tried to list rooms with date not in range, min : 1 max : %d", c.Username, Config.DayNumber)
		c.Msg("Date need to be in range. usage: /rooms DATE (NUMBER_NIGHTS min: 1 max: " + strconv.Itoa(Config.DayNumber) + ")")
		return
	}

	if err != nil {
		return
	}

	// show the list of the rooms and their availability
	if len(args) == 2 {
		var rooms []string

		if Config.Debug {
			log.Printf("Client (%s) in rooms function ask for Critical Section", c.Conn.RemoteAddr().String())
			//time.Sleep(5 * time.Second)
		}
		h.Request(c)
		if Config.Debug {
			log.Printf("Client (%s) in rooms function is in Critical Section", c.Conn.RemoteAddr().String())
			//time.Sleep(5 * time.Second)
		}
		// SECTION CRITIQUE
		for roomNumber := 0; roomNumber < Config.RoomsNumber; roomNumber++ {
			if h.Rooms[roomNumber][dateReservation-1] == "" {
				rooms = append(rooms, fmt.Sprintf("room number: %d is free", roomNumber+1))
			} else if h.Rooms[roomNumber][dateReservation] == c.Username { // TODO : N'a pas l'air de fonctionner
				rooms = append(rooms, fmt.Sprintf("room number: %d is occupied by you", roomNumber+1))
			} else {
				rooms = append(rooms, fmt.Sprintf("room number: %d is occupied", roomNumber+1))
			}
		}
		c.Msg(fmt.Sprintf("Rooms: %s", strings.Join(rooms, ", ")))
		// FIN SECTION CRITIQUE
		if Config.Debug {
			log.Printf("Client (%s) in rooms function finish with Critical Section", c.Username)
			//time.Sleep(5 * time.Second)
		}
		h.Release()
		if Config.Debug {
			log.Printf("Client (%s) in rooms function is out of Critical Section", c.Username)
			//time.Sleep(5 * time.Second)
		}
	} else {
		// show the first room ready to be used in the lap of time desired

		numberNightsRaw := args[2]
		numberNights, err := strconv.Atoi(numberNightsRaw)
		if err != nil {
			return
		}

		if dateReservation+numberNights > Config.DayNumber {
			log.Printf("Client (%s) tried to list rooms with date not in range, min : 1 max : %d", c.Username, Config.DayNumber)
			c.Msg("Date + Number of night need to be in range. usage: /rooms DATE (NUMBER_NIGHTS min: 1 max: " + strconv.Itoa(Config.DayNumber) + ")")
			return
		}
		if Config.Debug {
			log.Printf("Client (%s) in rooms function ask for Critical Section", c.Conn.RemoteAddr().String())
			//time.Sleep(5 * time.Second)
		}
		h.Request(c)
		if Config.Debug {
			log.Printf("Client (%s) in rooms function is in Critical Section", c.Conn.RemoteAddr().String())
			//time.Sleep(5 * time.Second)
		}
		// SECTION CRITIQUE
		for i := 0; i < Config.RoomsNumber; i++ {
			ok := true
			for j := dateReservation - 1; j < dateReservation+numberNights-1; j++ {
				if h.Rooms[i][j] != "" {
					ok = false
					break
				}
			}
			if ok {
				c.Msg(fmt.Sprintf("You can use this room : %d", i+1))
				return
			}
		}
		c.Msg(fmt.Sprintf("Sorry, no rooms available"))
		// FIN SECTION CRITIQUE
		if Config.Debug {
			log.Printf("Client (%s) in rooms function finish with Critical Section", c.Username)
			//time.Sleep(5 * time.Second)
		}
		h.Release()
		if Config.Debug {
			log.Printf("Client (%s) in rooms function is out of Critical Section", c.Username)
			//time.Sleep(5 * time.Second)
		}
	}
	log.Printf("Client (%s) successfully list rooms", c.Username)
}

// quit close the tcp connection between the Client and server
func (h *Hotel) quit(c *Network.Client) {
	log.Printf("Client has left the chat: %s", c.Username)
	c.Msg("Good bye!")
	c.Conn.Close()
}

// UpdateUsername called by Mutex to update clients array
func (h *Hotel) UpdateUsername(username string) {
	h.clients = append(h.clients, username)
}

// UpdateRooms called by Mutex to updates rooms map
func (h *Hotel) UpdateRooms(arguments []string) {

	roomNameRaw := arguments[0]
	dateReservationRaw := arguments[1]
	numberNightsRaw := arguments[2]
	username := arguments[3]

	roomName, err := strconv.Atoi(roomNameRaw)
	if err != nil {
		return
	}

	dateReservation, err := strconv.Atoi(dateReservationRaw)
	if err != nil {
		return
	}

	numberNights, err := strconv.Atoi(numberNightsRaw)
	if err != nil {
		return
	}

	//demander section critique
	h.Raymond.handleRequest()
	ok := <-h.AccessCS
	if ok {
		for i := dateReservation; i < dateReservation+numberNights; i++ {
			h.Rooms[roomName-1][i-1] = username
		}
	}
	//relacher section critique
	h.Raymond.handleRelease()
}

func (h *Hotel) Request(c *Network.Client) {
	h.MsgSC <- "req"
	ok := <-h.AccessCS
	if !ok {
		log.Printf("Client (%s) tried to enter in SC without success", c.Conn.RemoteAddr().String())
		c.Msg("Critical section was not available")
		return
	}
}

func (h *Hotel) Release() {
	h.MsgSC <- "rel"
}
