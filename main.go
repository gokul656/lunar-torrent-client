package main

import (
	"flag"
	"fmt"
	"os"
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
}
