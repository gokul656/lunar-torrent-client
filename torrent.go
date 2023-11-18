package main

import (
	"net/url"
	"strconv"
)

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PicesHash   [][20]byte
	PicesLength int
	Length      int
	Name        string
}

func (t *TorrentFile) toTorrentFile() (*TorrentFile, error) {
	return nil, nil
}

func (t *TorrentFile) buildTracekerURL(peerID [20]byte, port uint16) (string, error) {
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
