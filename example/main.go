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

	cfg := &tls.Config{
		MinVersion:   tls.VersionTLS13,
		Certificates: []tls.Certificate{keyPair},
	}

	m := tlsmux.Muxer{}
	m.Handle("foo.localhost", tlsmux.TLSHandlerFunc(cfg, func(conn net.Conn) {
		_, _ = conn.Write([]byte("foo"))
	}))
	m.Handle("bar.localhost", tlsmux.TLSHandlerFunc(cfg, func(conn net.Conn) {
		_, _ = conn.Write([]byte("bar"))
	}))

	l, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			panic(err)
		}

		go m.Serve(conn)
	}
}
