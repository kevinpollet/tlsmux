package tlsmux

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMuxer_Handle(t *testing.T) {
	m := Muxer{}
	m.Handle("Foo", HandlerFunc(func(conn net.Conn) {}))

	assert.Len(t, m.hs, 1)
	assert.Contains(t, m.hs, "foo")
}
