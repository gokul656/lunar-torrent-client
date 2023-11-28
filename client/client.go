package client

import (
	"bytes"
	"fmt"
	"net"

	"github.com/gokul656/lunar-torrent/peer"
)

type Client struct {
	Conn       net.Conn
	State      string
	PeerID     string
	Downloaded uint
	Uploaded   uint
	Pending    uint
	Peers      []peer.Peer
}

func (c *Client) InitateHandshake() error {
	conn := c.Conn

	handshakeBuf := []byte{}
	_, err := conn.Write(handshakeBuf)
	if err != nil {
		return err
	}

	responseBuf := bytes.Buffer{}
	_, err = conn.Read(responseBuf.Bytes())
	if err != nil {
		return err
	}

	fmt.Println(responseBuf)
	return nil
}
