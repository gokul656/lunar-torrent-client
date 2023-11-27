package torrent_file

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	"github.com/jackpal/bencode-go"
)

type BencodeTorrent struct {
	Announce  string      `bencode:"announce" json:"announce"`
	Comment   string      `bencode:"comment" json:"comment"`
	Info      BencodeInfo `bencode:"info" json:"info"`
	CreatedBy string      `bencode:"created by" json:"created by"`
	URLList   []string    `bencode:"url-list" json:"-"`
}

type BencodeInfo struct {
	Name        string `bencode:"name" json:"name"`
	Length      int    `bencode:"length" json:"length"`
	PieceLength int    `bencode:"piece length" json:"piece length"`
	Pieces      string `bencode:"pieces" json:"-"`
}

type BencodeFile struct {
	Length int    `bencode:"length" json:"length,omitempty"`
	Path   string `bencode:"path" json:"path,omitempty"`
	Pieces string `bencode:"pieces" json:"pieces,omitempty"`
}

// BencodeTorrent methods

func (bt *BencodeTorrent) ToTorrentFile() (*TorrentFile, error) {
	infoHash, err := bt.Info.hash()
	if err != nil {
		return nil, err
	}

	piecesHash, err := bt.Info.SplitHashes()
	if err != nil {
		return nil, err
	}

	torrentFile := &TorrentFile{
		Announce:    bt.Announce,
		InfoHash:    infoHash,
		PicesLength: bt.Info.PieceLength,
		Length:      bt.Info.Length,
		Name:        bt.Info.Name,
		PicesHash:   piecesHash,
	}

	return torrentFile, nil
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

func (b BencodeTorrent) ToJson() []byte {
	marhshalled, err := json.Marshal(b)
	if err != nil {
		fmt.Println(err)
	}

	return marhshalled
}

// BencodeInfo methods

func (b BencodeInfo) Print() {
	hasher := sha1.New()
	hasher.Write([]byte(b.Pieces))
	checkSum := hasher.Sum(nil)
	pieceHash := hex.EncodeToString(checkSum)

	fmt.Println("\tName          :", b.Name)
	fmt.Println("\tPices         :", pieceHash)
	fmt.Printf("\tPieces Size   : %v KiB\n", b.PieceLength/1024)
	fmt.Printf("\tTotal Size    : %v MiB\n", b.Length/1024/1024)
}

func (bi *BencodeInfo) SplitHashes() ([][20]byte, error) {
	hashLength := 20

	hashBuffer := []byte(bi.Pieces)
	hashCount := len(hashBuffer) / hashLength
	hashes := make([][20]byte, hashCount)

	for i := 0; i < hashLength; i++ {
		// splitting into 20 bytes
		slicedHash := hashBuffer[i*hashLength : (i+1)*hashLength]
		copy(hashes[i][:], slicedHash)
	}

	return hashes, nil
}

// Write implements the io.Writer interface for BencodeInfo
func (bi *BencodeInfo) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (bi *BencodeInfo) hash() ([20]byte, error) {
	byteBuffer := &bytes.Buffer{}
	err := bencode.Marshal(byteBuffer, *bi)
	if err != nil {
		return [20]byte{}, err
	}

	return sha1.Sum(byteBuffer.Bytes()), nil
}

func parseFile(file string) (*BencodeTorrent, error) {
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
