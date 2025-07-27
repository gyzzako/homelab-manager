package url

import (
	"errors"
	"homelab-manager/internal/hosts/providers"
	"io"
	"net/http"

	"gopkg.in/yaml.v3"
)

type URLProvider struct {
	URL   string
	Token string
}

func (p *URLProvider) GetHostEntries() ([]providers.HostEntry, error) {
	req, err := http.NewRequest("GET", p.URL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+p.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, errors.New("failed to fetch URL: " + resp.Status)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var entries []providers.HostEntry
	if err := yaml.Unmarshal(data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
