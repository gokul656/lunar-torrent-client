package torrent_file

func Read(path string) (*TorrentFile, error) {
	bencodeData, err := parseFile(path)
	if err != nil {
		return nil, nil
	}

	tf, err := bencodeData.ToTorrentFile()
	if err != nil {
		return nil, nil
	}

	return tf, nil
}
