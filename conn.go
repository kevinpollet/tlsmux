package tlsmux

import "net"

type conn struct {
	net.Conn

	peeked []byte
}

// Read reads data from the peeked bytes first and after from the wrapped connection.
func (c *conn) Read(b []byte) (int, error) {
	if len(c.peeked) == 0 {
		return c.Conn.Read(b)
	}

	n := copy(b, c.peeked)
	c.peeked = c.peeked[n:]

	return n, nil
}
