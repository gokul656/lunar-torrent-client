package main

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jackpal/bencode-go"
)

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

func (t *TorrentFile) getPeerList(peerID [20]byte, port uint16) ([]Peer, error) {
	requestUrl, err := t.BuildTracekerURL(peerID, port)
	if err != nil {
		return nil, err
	}

	c := &http.Client{Timeout: 15 * time.Second}
	resp, err := c.Get(requestUrl)
	if err != nil {
		return nil, err
	}

	trackerResponse := &BencodeTrackerResp{}
	err = bencode.Unmarshal(resp.Body, trackerResponse)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return UnmarshallPeers([]byte(trackerResponse.Peers))
}
