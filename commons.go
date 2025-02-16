package main

import "net/url"

var FallbackTrackers = []string{
	"udp://tracker.opentrackr.org:1337/announce",
	"udp://tracker.leechers-paradise.org:6969/announce",
	"udp://open.stealth.si:80/announce",
	"udp://tracker.torrent.eu.org:451/announce",
	"udp://tracker.coppersurfer.tk:6969/announce",
}

func ParseURL(requestURL string) (string, string, error) {
	trackerURL, err := url.Parse(requestURL)
	if err != nil {
		return "", "", err
	}

	return trackerURL.Hostname(), trackerURL.Port(), nil
}
