package peer

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
	if len(byteBuffer)%peerSize != 0 {
		return nil, errors.New("malformed peer list")
	}

	peers := make([]Peer, peerCount)
	for i := 0; i < peerCount; i++ {
		offset := i * peerSize
		peers[i] = Peer{
			IP:   net.IP(byteBuffer[offset : offset+4]),
			Port: binary.BigEndian.Uint16([]byte(byteBuffer[offset+4 : offset+6])),
		}
	}

	return peers, nil
}

// overriden method
func (p Peer) String() string {
	return net.JoinHostPort(p.IP.String(), strconv.Itoa(int(p.Port)))
}
