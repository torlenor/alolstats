package utils

import (
	"fmt"
	"regexp"
	"strconv"
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

// SplitNumericVersion splits a version string into major, minor, patch version,
// e.g., 9.1.3 or 5.74.11 are valid versions.
func SplitNumericVersion(version string) ([]uint32, error) {
	versionRegex, _ := regexp.Compile(`^(\d+)\.(\d+)\.(\d+)$`)

	if !versionRegex.MatchString(version) {
		return nil, fmt.Errorf("%s is not a valid version string", version)
	}

	versionStrings := versionRegex.FindAllStringSubmatch(version, 3)
	if len(versionStrings) != 1 && len(versionStrings[0]) != 4 {
		return nil, fmt.Errorf("Something bad happened parsing version string %s", version)
	}

	versions := []uint32{}
	for _, str := range versionStrings[0][1:] {
		i, err := strconv.Atoi(str)
		if err != nil || i < 0 {
			return nil, fmt.Errorf("Could not convert %s to a valid unsigned version integer", str)
		}
		versions = append(versions, uint32(i))
	}

	return versions, nil
}

// SplitNumericMatchVersion splits a match game version string into major, minor, patch, additional version,
// e.g., 9.1.3.123 or 5.74.11.432 are valid versions.
func SplitNumericMatchVersion(version string) ([]uint32, error) {
	versionRegex, _ := regexp.Compile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)$`)

	if !versionRegex.MatchString(version) {
		return nil, fmt.Errorf("%s is not a valid version string", version)
	}

	versionStrings := versionRegex.FindAllStringSubmatch(version, 4)
	if len(versionStrings) != 1 && len(versionStrings[0]) != 5 {
		return nil, fmt.Errorf("Something bad happened parsing version string %s", version)
	}

	versions := []uint32{}
	for _, str := range versionStrings[0][1:] {
		i, err := strconv.Atoi(str)
		if err != nil || i < 0 {
			return nil, fmt.Errorf("Could not convert %s to a valid unsigned version integer", str)
		}
		versions = append(versions, uint32(i))
	}

	return versions, nil
}
