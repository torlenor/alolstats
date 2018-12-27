package utils

import (
	"fmt"
	"strings"
)

// GenerateStatusResponse generates a Json string with statusCode and statusMessage specified
// This can be used to have consistent error responses
func GenerateStatusResponse(statusCode uint16, statusText string) string {
	return fmt.Sprintf(`{"status": { "status_code": %d, "message": "%s" } }`, statusCode, statusText)
}

// CleanUpSummonerName removes all invalid chars such that a consistent Riot API request can be made for a summoner and such
// that this summoner can be found again
func CleanUpSummonerName(summonerName string) string {
	summonerName = strings.ToLower(summonerName)
	summonerName = strings.Replace(summonerName, " ", "", -1)
	summonerName = strings.Replace(summonerName, "_", "", -1)
	summonerName = strings.Replace(summonerName, "-", "", -1)
	return summonerName
}
