package playback

type URIser interface {
	URIs() []string
}

type Playlister interface {
	SetPlaylist(URIser) error
}

type Player interface {
	Play()
}

type Pauser interface {
	Pause()
	UnPause()
}

type Randomer interface {
	SetRandom()
	UnsetRandom()
}

type Statuser interface {
	Status() (map[string]string, error)
}

type Quiter interface {
	Quit() error
}

type Playbacker interface {
	Playlister
	Player
	Pauser
	Randomer
	Statuser
	Quiter
}
