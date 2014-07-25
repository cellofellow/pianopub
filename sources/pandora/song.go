package pandora

import (
	"github.com/cellofellow/gopiano"
)

type Song struct {
	URI    string
	Artist string
	Album  string
	Name   string
	ArtURL string
	Token  string
	client *gopiano.Client
}

type Playlist []Song

func (p Playlist) URIs() []string {
	ret := make([]string, len(p))
	for i, s := range p {
		ret[i] = s.URI
	}
	return ret
}
