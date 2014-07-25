package pandora

import (
	"errors"

	"github.com/cellofellow/gopiano"
)

type Pandora struct {
	client   *gopiano.Client
	config   map[string]string
	Stations Stations
}

func NewPandora(config map[string]string) (*Pandora, error) {
	var err error
	client, err := gopiano.NewClient(gopiano.AndroidClient)
	if err != nil {
		return nil, err
	}
	if _, ok := config["username"]; !ok {
		return nil, errors.New("username not in config")
	}
	if _, ok := config["password"]; !ok {
		return nil, errors.New("password not in config")
	}
	_, err = client.AuthPartnerLogin()
	if err != nil {
		return nil, err
	}

	_, err = client.AuthUserLogin(config["username"], config["password"])
	if err != nil {
		return nil, err
	}

	return &Pandora{
		client: client,
		config: config,
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
