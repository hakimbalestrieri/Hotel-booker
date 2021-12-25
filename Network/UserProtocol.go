/**
File: UserProtocol.go
Authors: Hakim Balestrieri & Alexandre Mottier
Date: 22.10.2021
*/
package Network

type commandID int

const (
	CMD_USERNAME commandID = iota
	CMD_RESERVE
	CMD_ROOMS
	CMD_QUIT
)

// UserProtocol represent command send from tcp Client
type UserProtocol struct {
	Id commandID
	// Client using the command
	Client *Client
	// Args argument of the command. for exemple: /rooms 1, 1 is an argument
	Args []string
}