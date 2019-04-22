package riotclient

import "time"

// RuneReforged contains the information for one Rune Reforged
type RuneReforged struct {
	ID        int    `json:"id"`
	Key       string `json:"key"`
	Icon      string `json:"icon"`
	Name      string `json:"name"`
	ShortDesc string `json:"shortDesc"`
	LongDesc  string `json:"longDesc"`
}

// RunesReforgedList is used to pass around a list of Runes Reforged
type RunesReforgedList []struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Icon  string `json:"icon"`
	Name  string `json:"name"`
	Slots []struct {
		Runes []RuneReforged `json:"runes"`
	} `json:"slots"`

	Timestamp time.Time `json:"timestamp"`
}
