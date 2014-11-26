package pandora

import (
	"github.com/cellofellow/gopiano"
)

type Station struct {
	Name       string
	ArtURL     string
	Id         string
	Token      string
	IsQuickMix bool
	client     *gopiano.Client
}

func (s *Station) Playlist() (Playlist, error) {
	response, err := s.client.StationGetPlaylist(s.Token)
	if err != nil {
		return nil, err
	}

	songs := make(Playlist, len(response.Result.Items))
	for i, item := range response.Result.Items {
		songs[i] = Song{
			URI:    item.AudioURLMap["medium"].AudioURL,
			Artist: item.ArtistName,
			Album:  item.AlbumName,
			ArtURL: item.AlbumArtURL,
			Token:  item.TrackToken,
			client: s.client,
		}
	}
	return songs, nil
}

type Stations []Station

// Sort interface methods.

func (s Stations) Len() int {
	return len(s)
}

func (s Stations) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Stations) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
