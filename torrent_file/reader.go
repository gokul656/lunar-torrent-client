package torrent_file

func Read(path string) (*TorrentFile, error) {
	bencodeData, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	bencodeData.Print()

	tf, err := bencodeData.ToTorrentFile()
	if err != nil {
		return nil, err
	}

	return tf, nil
}
