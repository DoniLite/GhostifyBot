package services

import (
    "fmt"

    tr "github.com/anacrolix/torrent"
)

// Downloading a torrent file from source based on the provided URL
func DownloadTorrentFile(url, destinationDir string) (string, error) {
    return "", fmt.Errorf("not implemented function")
}

// Downloading a torrent file specified in a filepath directory.
func DownloadFromTorrentFile(torrentFilePath, downloadDir string) error {
    // creating new torrent
    clientConfig := tr.NewDefaultClientConfig()
    client, err := tr.NewClient(clientConfig)
    if err != nil {
        return fmt.Errorf("error during the torrent client creation : %v", err)
    }
    defer client.Close()

    torrent, err := client.AddTorrentFromFile(torrentFilePath)
    if err != nil {
        return fmt.Errorf("error during the torrent adding : %v", err)
    }

	clientConfig.DataDir = downloadDir

    // Waiting to get the metadata
    <-torrent.GotInfo()

    torrent.DownloadAll()
    
    client.WaitAll()

    return nil
}

// Download torrent file specified by the magnet link.
func DownloadFromMagnetLink(magnetLink, downloadDir string) error {
    clientConfig := tr.NewDefaultClientConfig()
    client, err := tr.NewClient(clientConfig)
    if err != nil {
        return fmt.Errorf("error during the torrent client creation : %v", err)
    }
    defer client.Close()

    // Adding the magnet link
    torrent, err := client.AddMagnet(magnetLink)
    if err != nil {
        return fmt.Errorf("error during the magnet link adding : %v", err)
    }

	clientConfig.DataDir = downloadDir

    <-torrent.GotInfo()

    torrent.DownloadAll()

    
    client.WaitAll()

    return nil
}


// Util func to check if the provided link is a magnet link
func IsMagnet(link string) bool {
	return len(link) > 8 && link[:8] == "magnet:?"
}