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
	fmt.Println("Target :", *targetFile)

	// Parse torrent file
	bencode, err := parseBencodeFile(*targetFile)
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	bencode.Print()
}
