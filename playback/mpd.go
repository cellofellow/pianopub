package playback

import (
	"time"

	"github.com/turbowookie/gompd/mpd"
)

type mpdplayer struct {
	client *mpd.Client
	quit   chan struct{}
}

func NewMpdPlaybacker(host string) (Playbacker, error) {
	client, err := mpd.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	player := &mpdplayer{
		client: client,
		quit:   make(chan struct{}, 1),
	}

	player.startManager()

	return player, nil
}

func (m *mpdplayer) startManager() {
	go func() {
		for {
			select {
			case <-m.quit:
				m.client.Close()
				return
			case <-time.Tick(30 * time.Second):
				m.client.Ping()
			}
		}
	}()
}

func (m *mpdplayer) SetPlaylist(item URIser) error {
	m.client.Clear()

	for _, uri := range item.URIs() {
		err := m.client.Add(uri)
		if err != nil {
			m.client.Clear()
			return err
		}
	}

	return nil

}

func (m *mpdplayer) Play() {
	m.client.Play(0)
}

func (m *mpdplayer) Pause() {
	m.client.Pause(true)
}

func (m *mpdplayer) UnPause() {
	m.client.Pause(false)
}

func (m *mpdplayer) SetRandom() {
	m.client.Random(true)
}

func (m *mpdplayer) UnsetRandom() {
	m.client.Random(false)
}

func (m *mpdplayer) Status() (map[string]string, error) {
	return m.client.Status()
}

func (m *mpdplayer) Quit() error {
	m.quit <- struct{}{}
	return nil
}
