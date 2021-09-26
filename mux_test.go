package tlsmux

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMuxer_Handle(t *testing.T) {
	m := Mux{}

	m.Handle("Foo", HandlerFunc(func(_ net.Conn) error {
		return nil
	}))

	assert.Len(t, m.hs, 1)
	assert.Contains(t, m.hs, "foo")
}
