package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
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

	// Parse torrent file
	bencode, err := parseBencodeFile(*targetFile)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	var peerId [20]byte
	copy(peerId[:], []byte("gokul656"))

	bencode.Print()

	torrentFile, _ := bencode.ToTorrentFile()
	peerList, err := torrentFile.getPeerList(peerId, 2323)
	if err != nil {
		panic(err)
	}

	fmt.Println("# Peers")
	fmt.Printf("\tCount         : %v\n", len(peerList))

	peerID := [20]byte{}
	rand.Read(peerID[:])

	var wg sync.WaitGroup
	for _, peer := range peerList {
		wg.Add(1)
		go func(peer Peer, wg *sync.WaitGroup) {
			defer wg.Done()
			conn, err := net.DialTimeout("tcp", peer.String(), time.Second*3)
			if err != nil {
				return
			}

			_, err = InitiateHandshake(conn, peerID, torrentFile.InfoHash)
			if err != nil {
				return
			}

			log.Println("Handshake success with", peer.IP.String())
			defer conn.Close()
		}(peer, &wg)
	}

	wg.Wait()
}
