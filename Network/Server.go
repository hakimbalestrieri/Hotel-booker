package Network

import "net"

type Server struct {
	id        int
	hasAccess bool
	conn      *net.Conn
}
