package main

import (
	"net"

	"github.com/kevinpollet/tlsmux"
)

func main() {
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	m := tlsmux.Mux{}
	m.Handle("foo.localhost", tlsmux.HandlerFunc(func(conn net.Conn) {

		_, _ = conn.Write([]byte("foo"))
	}))

	m.Handle("bar.localhost", tlsmux.HandlerFunc(func(conn net.Conn) {

		_, _ = conn.Write([]byte("bar"))
	}))

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go m.Serve(conn)
	}
}
