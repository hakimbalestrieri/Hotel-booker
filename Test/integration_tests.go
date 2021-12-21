package main

import (
	"PRR-Labo3-Balestrieri/Config"
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

func main() {
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(Config.ServerPorts[0]))
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	conn2, err := net.Dial("tcp", "localhost:"+strconv.Itoa(Config.ServerPorts[2]))
	if err != nil {
		log.Fatal(err)
	}

	defer conn2.Close()

	serverReader := bufio.NewReader(conn)
	serverReader2 := bufio.NewReader(conn2)

	testHotelCommand(1, conn, serverReader, "/username", "username is required. usage: /username NAME")

	testHotelCommand(2, conn, serverReader, "/username alex", "Welcome in the hotel, alex")

	testHotelCommand(3, conn, serverReader, "/reserve 1", "room name, date of arrival and number of nights are required. usage: /reserve ROOM_NUMBER DATE NUMBER_NIGHTS")

	testHotelCommand(4, conn, serverReader, "/reserve 1 1 1", "Room number: 1 is now reserved")

	testHotelCommand(5, conn2, serverReader2, "/username alex", "username already exist. usage: /username NAME")

	testHotelCommand(6, conn2, serverReader2, "/username hakim", "Welcome in the hotel, hakim")

	testHotelCommand(7, conn2, serverReader2, "/reserve 1 1 1", "The room is already occupied.")

	testHotelCommand(8, conn2, serverReader2, "/rooms 1", "Rooms: room number: 1 is occupied, room number: 2 is free, room number: 3 is free, room number: 4 is free, room number: 5 is free, room number: 6 is free, room number: 7 is free, room number: 8 is free, room number: 9 is free, room number: 10 is free, room number: 11 is free, room number: 12 is free, room number: 13 is free, room number: 14 is free, room number: 15 is free, room number: 16 is free, room number: 17 is free, room number: 18 is free, room number: 19 is free, room number: 20 is free, room number: 21 is free, room number: 22 is free, room number: 23 is free, room number: 24 is free, room number: 25 is free, room number: 26 is free, room number: 27 is free, room number: 28 is free, room number: 29 is free, room number: 30 is free")

	testHotelCommand(9, conn2, serverReader2, "/rooms 1 2", "You can use this room : 2")

	testHotelCommand(10, conn, serverReader, "/reserve 2 5 10", "Room number: 2 is now reserved")

	testHotelCommand(11, conn, serverReader, "/reserve 500 1 1", "room number was not in range : min 1, max 30")

	testHotelCommand(12, conn, serverReader, "/reserve 4 400 45", "duration was not in range : min 1, max 31")

	testHotelCommand(13, conn, serverReader, "/rooms 45", "Date need to be in range. usage: /rooms DATE (NUMBER_NIGHTS min: 1 max: 31)")
}

func testHotelCommand(testNumber int, conn net.Conn, reader *bufio.Reader, command string, shouldReturn string) {
	fmt.Println(fmt.Sprintf("TEST - %d : %s with user connexion : %s", testNumber, command, conn.RemoteAddr().String()))
	conn.Write([]byte(command + "\n"))

	serverResponse, _ := reader.ReadString('\n')
	if serverResponse != "> "+shouldReturn+"\n" {
		fmt.Println(fmt.Sprintf("FAILED: TEST - %d! Should return : %s, but returned : %s", testNumber, shouldReturn, serverResponse))
		return
	}
	fmt.Println(fmt.Sprintf("SUCCESS: TEST - %d! return : %s", testNumber, serverResponse))
}
