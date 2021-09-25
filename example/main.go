package main

import (
	"crypto/tls"
	"net"

	"github.com/kevinpollet/tlsmux"
)

func main() {
	keyPair, err := tls.LoadX509KeyPair("./example/certs/_.localhost/cert.pem", "./example/certs/_.localhost/key.pem")
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{keyPair}}

	mux := tlsmux.Mux{}

	mux.Handle("foo.localhost", tlsmux.TLSHandlerFunc(cfg, tlsmux.HandlerFunc(func(conn net.Conn) {
		_, _ = conn.Write([]byte("foo"))
	})))

	mux.Handle("bar.localhost", tlsmux.TLSHandlerFunc(cfg, tlsmux.HandlerFunc(func(conn net.Conn) {
		_, _ = conn.Write([]byte("bar"))
	})))

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
