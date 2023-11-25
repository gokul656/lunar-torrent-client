package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Handshake struct {
	Pstr     string
	InfoHash [20]byte
	PeerID   [20]byte
}

func NewHandshake(peerID, infoHash [20]byte) *Handshake {
	return &Handshake{
		Pstr:     "BitTorrent protocol", // case sensitive & should never be changed
		InfoHash: infoHash,
		PeerID:   peerID,
	}
}

func InitiateHandshake(conn net.Conn, peerID, infoHash [20]byte) (*Handshake, error) {
	conn.SetDeadline(time.Now().Add(time.Second * 3))
	defer conn.SetDeadline(time.Time{}) // Disable the deadline

	hs := NewHandshake(peerID, infoHash)
	_, err := conn.Write(hs.Serialize())
	if err != nil {
		return nil, fmt.Errorf("unable to write to connection %v", err.Error())
	}

	_, err = Read(conn)
	if err != nil {
		return nil, fmt.Errorf("unable to read from connection %v", err.Error())
	}

	return hs, nil
}

func (h *Handshake) Serialize() []byte {
	buf := make([]byte, len(h.Pstr)+49)
	buf[0] = byte(len(h.Pstr))
	curr := 1
	curr += copy(buf[curr:], h.Pstr)
	curr += copy(buf[curr:], make([]byte, 8)) // 8 reserved bytes
	curr += copy(buf[curr:], h.InfoHash[:])
	curr += copy(buf[curr:], h.PeerID[:])
	return buf
}

func Read(r io.Reader) (*Handshake, error) {
	lengthBuf := make([]byte, 1)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}

	pstrlen := int(lengthBuf[0])
	if pstrlen == 0 {
		err := fmt.Errorf("pstrlen cannot be 0")
		return nil, err
	}

	handshakeBuf := make([]byte, 48+pstrlen)
	_, err = io.ReadFull(r, handshakeBuf)
	if err != nil {
		log.Println("unable to read handshake buffer")
		return nil, err
	}

	var infoHash, peerID [20]byte

	copy(infoHash[:], handshakeBuf[pstrlen+8:pstrlen+8+20])
	copy(peerID[:], handshakeBuf[pstrlen+8+20:])

	handshake := &Handshake{
		Pstr:     string(handshakeBuf[0:pstrlen]),
		InfoHash: infoHash,
		PeerID:   peerID,
	}

	return handshake, nil
}
