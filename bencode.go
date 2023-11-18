package main

import (
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type BencodeTorrent struct {
	Announce  string      `bencode:"announce"`
	Comment   string      `bencode:"comment"`
	Info      BencodeInfo `bencode:"info"`
	CreatedBy string      `bencode:"created by"`
	URLList   []string    `bencode:"url-list"`
}

type BencodeInfo struct {
	Files       BencodeFile `bencode:"files"`
	Name        string      `bencode:"name"`
	Length      int         `bencode:"length"`
	PieceLength float64     `bencode:"piece length"`
	Pieces      string      `bencode:"pieces"`
}

type BencodeFile struct {
	Length float64 `bencode:"length"`
	Path   string  `bencode:"path"`
	Pieces any     `bencode:"pieces"`
}

func parseBencodeFile(file string) (*BencodeTorrent, error) {
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	bc := &BencodeTorrent{}
	err = bencode.Unmarshal(reader, bc)
	if err != nil {
		return nil, err
	}

	return bc, nil
}

func (b BencodeTorrent) Print() {
	fmt.Println("\n# Root")
	fmt.Println("\tAnnounce      :", b.Announce)
	fmt.Println("\tComment       :", b.Comment)
	fmt.Println("\tCreated by    :", b.CreatedBy)

	fmt.Println("# Info")
	b.Info.Print()

	if len(b.URLList) > 0 {
		fmt.Println("# URL List")
		for i, url := range b.URLList {
			fmt.Printf("\t%d             : %s\n", i, url)
		}
	}
}

func (b BencodeInfo) Print() {
	fmt.Println("\tName          :", b.Name)
	fmt.Printf("\tPieces Size   : %v KiB\n", b.PieceLength/1024)
	fmt.Printf("\tTotal Size    : %v MiB\n", b.Length/1024/1024)
}
