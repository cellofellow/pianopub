package pandora

import (
	"github.com/cellofellow/gopiano"
)

type Pandora struct {
	client             *gopiano.Client
	username, password string
	Stations           Stations
}

func NewPandora(username, password string) (*Pandora, error) {
	var err error
	client, err := gopiano.NewClient(gopiano.AndroidClient)
	if err != nil {
		return nil, err
	}
	_, err = client.AuthPartnerLogin()
	if err != nil {
		return nil, err
	}

	_, err = client.AuthUserLogin(username, password)
	if err != nil {
		return nil, err
	}

	return &Pandora{
		client:   client,
		username: username,
		password: password,
	}, nil
}

func (p *Pandora) FetchStations() error {
	response, err := p.client.UserGetStationList(true)
	if err != nil {
		return err
	}
	var stations Stations = make([]Station, len(response.Result.Stations))
	for i, st := range response.Result.Stations {
		stations[i] = Station{
			Name:       st.StationName,
			ArtURL:     st.ArtURL,
			Id:         st.StationID,
			Token:      st.StationToken,
			IsQuickMix: st.IsQuickMix,
			client:     p.client,
		}
	}
	p.Stations = stations
	return nil
}
