package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/gokul656/lunar-torrent/torrent_file"
)

var targetFile = flag.String("t", "", "Path to torrent file")

func init() {
	flag.Parse()
	if *targetFile == "" {
		panic("invalid torrent file")
	}

	_, err := os.Open(*targetFile)
	if os.IsExist(err) {
		panic("file not found")
	}
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("\nPanic due to,", r)
		}
	}()
	fmt.Println("Target :", *targetFile)

	tf, err := torrent_file.Read(*targetFile)
	if err != nil {
		panic(err)
	}

	tf.IntiateDownload("out/")

	// // Parse torrent file
	// bencode, err := parseBencodeFile(*targetFile)
	// if err != nil {
	// 	fmt.Printf("[ERROR] %s\n", err)
	// }

	// var peerId [20]byte
	// copy(peerId[:], []byte("gokul656"))

	// bencode.Print()

	// torrentFile, _ := bencode.ToTorrentFile()
	// peerList, err := torrentFile.getPeerList(peerId, 2323)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("# Peers")
	// fmt.Printf("\tCount         : %v\n", len(peerList))

	// peerID := [20]byte{}
	// copy(peerID[:], []byte("lunar-torrent-client"))

	// var wg sync.WaitGroup
	// for _, peer := range peerList {
	// 	wg.Add(1)
	// 	go func(peer Peer, wg *sync.WaitGroup) {
	// 		defer wg.Done()
	// 		conn, err := net.DialTimeout("tcp", peer.String(), time.Second*3)
	// 		if err != nil {
	// 			return
	// 		}
	// 		defer conn.Close()

	// 		_, err = InitiateHandshake(conn, peerID, torrentFile.InfoHash)
	// 		if err != nil {
	// 			return
	// 		}

	// 		log.Println("Handshake success with", peer.IP.String())

	// 		// wait for server to get ready (i.e) to receive unchoked msg
	// 		msg, err := ReadMessage(conn)
	// 		if err != nil {
	// 			return
	// 		}

	// 		fmt.Println("message id :", msg.ID)
	// 		fmt.Println("payload len:", len(msg.Payload))

	// 		for pieceIndex := range torrentFile.PicesHash {
	// 			if msg.Payload.HasPiece(pieceIndex) {
	// 				intrested := Message{
	// 					ID: MsgInterested,
	// 				}

	// 				conn.Write(intrested.Serialize())
	// 			}
	// 		}
	// 	}(peer, &wg)
	// }

	// wg.Wait()
}
