package riotclient

import (
	"time"
)

// ChampionInfo contains the basic info about the champion
type ChampionInfo struct {
	Attack     int `json:"attack"`
	Defense    int `json:"defense"`
	Magic      int `json:"magic"`
	Difficulty int `json:"difficulty"`
}

// ChampionImage contains information related to fetch the champion image
type ChampionImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

// ChampionStats contains extended information about the champion
type ChampionStats struct {
	Hp                   float64 `json:"hp"`
	Hpperlevel           float64 `json:"hpperlevel"`
	Mp                   float64 `json:"mp"`
	Mpperlevel           float64 `json:"mpperlevel"`
	Movespeed            float64 `json:"movespeed"`
	Armor                float64 `json:"armor"`
	Armorperlevel        float64 `json:"armorperlevel"`
	Spellblock           float64 `json:"spellblock"`
	Spellblockperlevel   float64 `json:"spellblockperlevel"`
	Attackrange          float64 `json:"attackrange"`
	Hpregen              float64 `json:"hpregen"`
	Hpregenperlevel      float64 `json:"hpregenperlevel"`
	Mpregen              float64 `json:"mpregen"`
	Mpregenperlevel      float64 `json:"mpregenperlevel"`
	Crit                 float64 `json:"crit"`
	Critperlevel         float64 `json:"critperlevel"`
	Attackdamage         float64 `json:"attackdamage"`
	Attackdamageperlevel float64 `json:"attackdamageperlevel"`
	Attackspeedoffset    float64 `json:"attackspeedoffset"`
	Attackspeedperlevel  float64 `json:"attackspeedperlevel"`
	Attackspeed          float64 `json:"attackspeed"`
}

// Champion stores champion data
type Champion struct {
	Version   string        `json:"version"`
	ID        string        `json:"id"`
	Key       string        `json:"key"`
	Name      string        `json:"name"`
	Title     string        `json:"title"`
	Blurb     string        `json:"blurb"`
	Info      ChampionInfo  `json:"info"`
	Image     ChampionImage `json:"image"`
	Tags      []string      `json:"tags"`
	Partype   string        `json:"partype"`
	Stats     ChampionStats `json:"stats"`
	Timestamp time.Time     `json:"timestamp"`
}

// ChampionList stores a list of Champion data
type ChampionList struct {
	Champions map[string]Champion `json:"champions"`
}
