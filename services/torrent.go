package services

import (
    "fmt"

    tr "github.com/anacrolix/torrent"
)

// Télécharge le fichier torrent depuis l'URL spécifiée et retourne le chemin du fichier téléchargé.
func DownloadTorrentFile(url, destinationDir string) (string, error) {
    // Implémentez le téléchargement du fichier torrent depuis l'URL.
    // Vous pouvez utiliser un package HTTP pour télécharger le fichier et l'enregistrer dans destinationDir.
    // Retournez le chemin complet du fichier téléchargé.
    return "", fmt.Errorf("fonction non implémentée")
}

// Télécharge les fichiers du torrent spécifié par le chemin du fichier torrent.
func DownloadFromTorrentFile(torrentFilePath, downloadDir string) error {
    // Crée un nouveau client torrent
    clientConfig := tr.NewDefaultClientConfig()
    client, err := tr.NewClient(clientConfig)
    if err != nil {
        return fmt.Errorf("erreur lors de la création du client torrent : %v", err)
    }
    defer client.Close()

    // Ajoute le torrent depuis le fichier
    torrent, err := client.AddTorrentFromFile(torrentFilePath)
    if err != nil {
        return fmt.Errorf("erreur lors de l'ajout du torrent : %v", err)
    }

	clientConfig.DataDir = downloadDir

    // Attendre que les métadonnées soient disponibles
    <-torrent.GotInfo()

    // Définir le répertoire de téléchargement
    torrent.DownloadAll()
    
    client.WaitAll()

    return nil
}

// Télécharge les fichiers du torrent spécifié par le magnet link.
func DownloadFromMagnetLink(magnetLink, downloadDir string) error {
    // Crée un nouveau client torrent
    clientConfig := tr.NewDefaultClientConfig()
    client, err := tr.NewClient(clientConfig)
    if err != nil {
        return fmt.Errorf("erreur lors de la création du client torrent : %v", err)
    }
    defer client.Close()

    // Ajoute le torrent depuis le magnet link
    torrent, err := client.AddMagnet(magnetLink)
    if err != nil {
        return fmt.Errorf("erreur lors de l'ajout du magnet link : %v", err)
    }

	clientConfig.DataDir = downloadDir

    // Attendre que les métadonnées soient disponibles
    <-torrent.GotInfo()

    // Définir le répertoire de téléchargement
    torrent.DownloadAll()

    
    client.WaitAll()

    return nil
}


// Fonction utilitaire pour vérifier si une chaîne est un lien magnet
func IsMagnet(link string) bool {
	return len(link) > 8 && link[:8] == "magnet:?"
}