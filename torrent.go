package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

type ProtocolHandler func(requestURL string) (*BencodeTrackerResp, error)

type BencodeTrackerResp struct {
	Interval int    `bencode:"interval"`
	Peers    string `bencode:"peers"`
}

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PicesHash   [][20]byte
	PicesLength int
	Length      int
	Name        string
}

func (t *TorrentFile) BuildTracekerURL(peerID [20]byte, port uint16) (string, error) {
	fmt.Printf("t.Announce: %v\n", t.Announce)
	trackerURL, err := url.Parse(t.Announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(t.InfoHash[:])},
		"peer_id":    []string{string(peerID[:])},
		"port":       []string{strconv.Itoa(int(port))},
		"uploaded":   []string{"0"},
		"downloaded": []string{"0"},
		"compact":    []string{"1"},
		"left":       []string{strconv.Itoa(t.Length)},
	}
	trackerURL.RawQuery = params.Encode()
	return trackerURL.String(), nil
}

func (t *TorrentFile) GetPeerList(peerID [20]byte, port uint16) ([]Peer, error) {
	handler := NewUDPHandler(t.Announce, port, peerID)
	peers, err := handler.GetPeerList()
	if err != nil {
		return nil, err
	}

	return peers, nil
}

// FIXME : Optimiza handlers
func httpHandler(requestURL string) (*BencodeTrackerResp, error) {
	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(requestURL)
	if err != nil {
		log.Println("unable to get peer list from URL", requestURL, err)
		return nil, err
	}

	trackerResponse := &BencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, trackerResponse)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return trackerResponse, nil
}

func udpHandler(requestURL string) (*BencodeTrackerResp, error) {
	conn, err := net.Dial("udp", requestURL)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	return nil, err
}
