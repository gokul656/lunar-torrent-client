package main

import (
	"encoding/binary"
	"errors"
	"net"
	"strconv"
)

type Peer struct {
	IP   net.IP
	Port uint16
}

func Unmarshall(byteBuffer []byte) ([]Peer, error) {
	const peerSize = 6 // 4 bytes for IP and 2 bytes for port
	peerCount := len(byteBuffer) / peerSize
	peers := make([]Peer, peerCount)

	if len(byteBuffer)%peerSize != 0 {
		return nil, errors.New("invalid peer")
	}

	for i := 0; i < peerCount; i++ {
		offset := i * peerSize
		peers[i] = Peer{
			IP:   net.IP(byteBuffer[offset : offset+4]),
			Port: binary.BigEndian.Uint16([]byte(byteBuffer[offset+4 : offset+6])),
		}
	}

	return peers, nil
}

func Address(peer Peer) string {
	return net.JoinHostPort(peer.IP.String(), strconv.Itoa(int(peer.Port)))
}
