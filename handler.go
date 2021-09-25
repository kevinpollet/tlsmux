package tlsmux

import "net"

type Handler interface {
	Serve(conn net.Conn)
}

type HandlerFunc func(conn net.Conn)

func (h HandlerFunc) Serve(conn net.Conn) {
	h(conn)
}
