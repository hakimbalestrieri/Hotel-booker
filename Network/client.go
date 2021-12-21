/**
File: Client.go
Authors: Hakim Balestrieri
Date: 22.10.2021
*/
package Network

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Client represent hotel's Client
type Client struct {
	// Conn tcp connection
	Conn     net.Conn
	Username string
	serverId int
	Commands chan UserProtocol
}

// ReadInput read Client inputs and affect a command if the protocol is respected
func (c *Client) ReadInput() {
	for {
		msg, err := bufio.NewReader(c.Conn).ReadString('\n')
		if err != nil {
			return
		}

		msg = strings.Trim(msg, "\r\n")

		args := strings.Split(msg, " ")
		cmd := strings.TrimSpace(args[0])
		switch cmd {
		case "/username":
			c.Commands <- UserProtocol{
				Id:     CMD_USERNAME,
				Client: c,
				Args:   args,
			}
		case "/reserve":
			c.Commands <- UserProtocol{
				Id:     CMD_RESERVE,
				Client: c,
				Args:   args,
			}
		case "/rooms":
			c.Commands <- UserProtocol{
				Id:     CMD_ROOMS,
				Client: c,
				Args:   args,
			}
		case "/quit":
			c.Commands <- UserProtocol{
				Id:     CMD_QUIT,
				Client: c,
			}
		default:
			c.err(fmt.Errorf("unknown command: %s", cmd))
		}
	}
}

// err write to tcp Client the error
func (c *Client) err(err error) {
	c.Conn.Write([]byte("err: " + err.Error() + "\n"))
}

// Msg write to tcp Client the message
func (c *Client) Msg(msg string) {
	c.Conn.Write([]byte("> " + msg + "\n"))
}
