package riotclientv4

import (
	"encoding/json"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
)

// Versions gets all versions from Data Dragon
func (c *RiotClientV4) Versions() (s riotclient.Versions, err error) {
	championsData, err := c.ddragon.GetLoLVersions()
	if err != nil {
		return nil, fmt.Errorf("Error getting versions from Data Dragon: %s", err)
	}

	versions := riotclient.Versions{}
	err = json.Unmarshal(championsData, &versions)
	if err != nil {
		return nil, err
	}

	return versions, nil
}
