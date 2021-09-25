package tlsmux

import "net"

type Conn struct {
	net.Conn

	peeked []byte
}

func (c *Conn) Read(b []byte) (int, error) {
	if len(c.peeked) == 0 {
		return c.Conn.Read(b)
	}

	n := copy(b, c.peeked)
	c.peeked = c.peeked[:n]

	return n, nil
}
