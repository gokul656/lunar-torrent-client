package client

import (
	"net"

	"github.com/gokul656/lunar-torrent/transport"
)

type ConnectionWrapper struct {
	Conn net.Conn
}

func (c *ConnectionWrapper) SendMessage(id transport.MessageID) {
	msg := transport.Message{
		ID: id,
	}

	c.Conn.Write(msg.Serialize())
}
