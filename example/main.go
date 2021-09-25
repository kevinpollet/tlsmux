package main

import (
	"net"

	"github.com/kevinpollet/tlsmux"
)

func main() {
	mux := tlsmux.Mux{}

	mux.Handle("foo.localhost", tlsmux.HandlerFunc(func(conn net.Conn) {
		_, _ = conn.Write([]byte("foo"))
	}))

	mux.Handle("bar.localhost", tlsmux.HandlerFunc(func(conn net.Conn) {
		_, _ = conn.Write([]byte("bar"))
	}))

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go mux.Serve(conn)
	}
}
