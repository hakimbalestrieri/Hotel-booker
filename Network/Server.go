package Network

import "net"

type Server struct {
	id   int
	conn *net.Conn
}
