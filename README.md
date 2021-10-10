# tlsmux

[![build](https://github.com/kevinpollet/tlsmux/actions/workflows/main.yml/badge.svg)](https://github.com/kevinpollet/tlsmux/actions)
[![GoDoc](https://godoc.org/github.com/kevinpollet/tlsmux?status.svg)](https://pkg.go.dev/github.com/kevinpollet/tlsmux)

Go package providing an implementation of a `net.Conn` multiplexer based on the TLS [SNI](https://www.cloudflare.com/learning/ssl/what-is-sni/) (Server Name Indication) sent by a client.

## Installation

Install using `go get github.com/kevinpollet/tlsmux`.

## Usage

### Mux

The `Mux` struct allows registering handlers which will be called when the muxer serve a `net.Conn` with a 
matching server name.

```go
mux := tlsmux.Mux{}

l, err := net.Listen("tcp", "127.0.0.1:8080")
if err != nil {
    panic(err)
}

if err := mux.Serve(l); err != nil {
    log.Fatal(err)
}
```

### Handler

The `Handler` interface is used to handle an incoming `net.Conn` without decrypting the underlying TLS communication (Pass Through).
Implementations are responsible for closing the connection.

```go
mux.Handle("server.name", tlsmux.HandlerFunc(func(conn net.Conn) error {
    defer conn.Close()

    dstConn, err := net.Dial("tcp", "127.0.0.1:443")
    if err != nil {
        return fmt.Errorf("dial: %w", err)
    }
    defer dstConn.Close()

    go func() { io.Copy(dstConn, conn) }()
    io.Copy(conn, dstConn)

    return nil
}))
```

### TLSHandler

The `TLSHandler` struct is a `Handler` implementation allowing to terminate the TLS connection with the configured `tls.Config`.
Thus, the `net.Conn` parameter of a `TLSHandler` if of type `tls.Conn`.  
Implementations are responsible for closing the connection.

```go
cfg := &tls.Config{
    MinVersion: tls.VersionTLS13,
    Certificates: []tls.Certificate{cert},
}

mux.Handle("foo.localhost", tlsmux.TLSHandlerFunc(cfg, func(conn net.Conn) error {
    defer conn.Close()

    _, err := io.WriteString(conn, "foo")
    return err
}))
```

### ProxyHandler

The `ProxyHandler` struct is a `Handler` implementation forwarding the connection bytes to the configured `Address`.

```go
mux.Handle("foo.localhost", tlsmux.ProxyHandler{Addr: "127.0.0.1:443"})
```

The `ProxyHandlerFunc` is an adapter allowing the use of a `ProxyHandler` as a `HandlerFunc`.

```go
cfg := &tls.Config{
    MinVersion: tls.VersionTLS13,
    Certificates: []tls.Certificate{cert},
}

mux.Handle("foo.localhost", tlsmux.TLSHandlerFunc(cfg, tlsmux.ProxyHandlerFunc("127.0.0.1:80"))
```

## License

[MIT](./LICENSE.md)
