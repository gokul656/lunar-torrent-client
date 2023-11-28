package torrent_file

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gokul656/lunar-torrent/client"
	"github.com/gokul656/lunar-torrent/peer"
	"github.com/jackpal/bencode-go"
)

type BencodeTrackerResponse struct {
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

func (tf *TorrentFile) getPeerList(peerID [20]byte, port uint16) ([]peer.Peer, error) {

	// to retrive peer list
	announceURL, err := tf.buildTracekerURL(peerID, 2032)
	if err != nil {
		return []peer.Peer{}, err
	}

	// announcing about our the prescence & retrieving peer list
	client := &http.Client{Timeout: time.Second * 10}
	responseBuf, err := client.Get(announceURL)
	if err != nil {
		return []peer.Peer{}, err
	}

	defer responseBuf.Body.Close()

	trackerResponse := &BencodeTrackerResponse{}
	err = bencode.Unmarshal(responseBuf.Body, trackerResponse)
	if err != nil {
		return []peer.Peer{}, err
	}

	return peer.Unmarshall([]byte(trackerResponse.Peers))
}

func (tf *TorrentFile) IntiateDownload(dst string) error {
	peerID := [20]byte{}
	copy(peerID[:], []byte("lunar-torrent-client"))

	peerList, err := tf.getPeerList(peerID, 2432)
	if err != nil {
		return err
	}

	torrent := &client.Torrent{
		PeerID:     peerID,
		Peers:      peerList,
		InfoHash:   tf.InfoHash,
		PiecesHash: tf.PicesHash,
	}

	return torrent.Download()
}
