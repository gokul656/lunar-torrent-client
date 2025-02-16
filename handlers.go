package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

type UDPHandler struct {
	infoHash [20]byte
	peerID   [20]byte
	port     uint16
	URL      string
}

func NewUDPHandler(URL string, port uint16, peerID [20]byte) *UDPHandler {
	return &UDPHandler{
		URL:    URL,
		port:   port,
		peerID: peerID,
	}
}

func (h *UDPHandler) BuildTrackerURL(announceURL string, peerID [20]byte) (string, error) {
	host, port, err := ParseURL(announceURL)
	if err != nil {
		return "", err
	}

	if port == "" {
		port = fmt.Sprint(h.port)
	}

	return net.JoinHostPort(host, port), nil
}

func generateTransactionID() uint32 {
	var id uint32
	binary.Read(rand.Reader, binary.BigEndian, &id)
	return id
}

func (h *UDPHandler) GetPeerList() ([]Peer, error) {
	url, err := h.BuildTrackerURL(h.URL, h.peerID)
	if err != nil {
		return nil, err
	}

	var conn net.Conn
	conn, err = connectTracker(url)
	if err != nil {
		for _, url := range FallbackTrackers {
			conn, err = connectTracker(url)
			if err != nil {
				fmt.Printf("err: %v\n", err)
				continue
			}

			// fmt.Println("Connection established with", url)

			break
		}
	}

	defer conn.Close()

	peers, _ := h.getPeers(conn)
	if len(peers) == 0 {
		return nil, errors.New("unable find peers")
	}

	return peers, nil
}

func connectTracker(url string) (net.Conn, error) {
	host, port, err := ParseURL(url)
	if err != nil {
		return nil, err
	}

	address := net.JoinHostPort(host, port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}

	return net.DialUDP("udp", nil, udpAddr)
}

func (h *UDPHandler) getPeers(conn net.Conn) ([]Peer, error) {

	transactionID := generateTransactionID()
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, uint64(0x41727101980)) // Magic constant
	binary.Write(&buf, binary.BigEndian, uint32(0))             // Connect action (0)
	binary.Write(&buf, binary.BigEndian, transactionID)         // Transaction ID

	_, err := conn.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending connect request: %v", err)
	}

	// Set timeout for response
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// Receive Connection Response
	resp := make([]byte, 16)
	_, err = conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("error receiving connect response: %v", err)
	}

	// Parse response
	action := binary.BigEndian.Uint32(resp[0:4])
	if action != 0 {
		return nil, fmt.Errorf("invalid connect response action")
	}

	connectionID := binary.BigEndian.Uint64(resp[8:16])

	buf.Reset()
	binary.Write(&buf, binary.BigEndian, connectionID)  // Connection ID
	binary.Write(&buf, binary.BigEndian, uint32(1))     // Announce action (1)
	binary.Write(&buf, binary.BigEndian, transactionID) // Transaction ID
	buf.Write(h.infoHash[:])                            // Info hash
	buf.Write(h.peerID[:])                              // Peer ID
	binary.Write(&buf, binary.BigEndian, uint64(0))     // Downloaded
	binary.Write(&buf, binary.BigEndian, uint64(0))     // Left
	binary.Write(&buf, binary.BigEndian, uint64(0))     // Uploaded
	binary.Write(&buf, binary.BigEndian, uint32(0))     // Event (None)
	binary.Write(&buf, binary.BigEndian, uint32(0))     // IP address (0)
	binary.Write(&buf, binary.BigEndian, uint32(0))     // Key (random)
	binary.Write(&buf, binary.BigEndian, int32(-1))     // Num want (-1 for all)
	binary.Write(&buf, binary.BigEndian, h.port)        // Port

	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error sending announce request: %v", err)
	}

	// Receive Announce Response
	resp = make([]byte, 1024)
	n, err := conn.Read(resp)
	if err != nil {
		return nil, fmt.Errorf("error receiving announce response: %v", err)
	}

	announceAction := binary.BigEndian.Uint32(resp[0:4])
	if announceAction != 1 {
		return nil, fmt.Errorf("invalid announce response action")
	}

	peerData := resp[20:n] // Peers start from byte 20

	return UnmarshallPeers(peerData)
}
