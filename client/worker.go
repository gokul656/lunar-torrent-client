package client

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/gokul656/lunar-torrent/peer"
	"github.com/gokul656/lunar-torrent/transport"
)

type Response struct {
	piece []byte
	err   error
}

type Torrent struct {
	Name       string
	Peers      []peer.Peer
	PeerID     [20]byte
	PiecesHash [][20]byte
	InfoHash   [20]byte
}

func (t *Torrent) Download() error {
	wg := &sync.WaitGroup{}
	responseChan := make(chan Response, len(t.Peers))
	pieceQueue := make(chan int)

	for pieceIndex, peer := range t.Peers {
		wg.Add(1)

		go t.connectWithPeer(peer, responseChan, pieceQueue, wg)
		go func(pieceIndex int) { pieceQueue <- pieceIndex }(pieceIndex)
	}

	go pieceCollector(responseChan)

	wg.Wait()
	close(responseChan)
	close(pieceQueue)
	return nil

}

func (t *Torrent) connectWithPeer(peer peer.Peer, responseChan chan Response, workerQueue chan int, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := net.DialTimeout("tcp", peer.String(), time.Second*1)
	if err != nil {
		responseChan <- createResponse(nil, err)
		return
	}

	defer conn.Close()

	_, err = transport.InitiateHandshake(conn, t.PeerID, t.InfoHash)
	if err != nil {
		responseChan <- createResponse(nil, err)
		return
	}

	message, err := transport.ReadMessage(conn)
	if err != nil {
		responseChan <- createResponse(nil, err)
		return
	}

	connectionWrapper := &ConnectionWrapper{Conn: conn}
	fmt.Println("Handshake success:", message.ID.Name())

	go downloadPiece(connectionWrapper, responseChan, workerQueue)
}

func downloadPiece(conn *ConnectionWrapper, responseChan chan Response, workerQueue chan int) {
	for pieceIndex := range workerQueue {
		fmt.Println("Downloading piece", pieceIndex)
		piece := fmt.Sprintf("%d", pieceIndex)

		conn.SendMessage(transport.MsgUnchoke)

		responseChan <- createResponse([]byte(piece), nil)
		continue
	}
}

func pieceCollector(responseChan <-chan Response) {
	for response := range responseChan {
		if response.err != nil {
			// fmt.Println("Error", response.err)
			continue
		} else {
			fmt.Println("Got piece!", string(response.piece))
		}
	}
}
func createResponse(piece []byte, err error) Response {
	return Response{
		piece: piece,
		err:   err,
	}
}
