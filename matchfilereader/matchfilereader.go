// Package matchfilereader is used to read the seed data provided from Riot
// See https://developer.riotgames.com/getting-started.html for details
package matchfilereader

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"git.abyle.org/hps/alolstats/riotclient"
)

// ReadMatchesFile reads the Riot provided matches file specified with filePath and tries to parse it
func ReadMatchesFile(filePath string) (*riotclient.Matches, error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		return nil, err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result riotclient.Matches
	json.Unmarshal([]byte(byteValue), &result)

	return &result, nil
}
