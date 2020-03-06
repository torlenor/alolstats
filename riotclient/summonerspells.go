package riotclient

import "time"

type SummonerSpellImage struct {
	Full   string `json:"full"`
	Sprite string `json:"sprite"`
	Group  string `json:"group"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	W      int    `json:"w"`
	H      int    `json:"h"`
}

// SummonerSpell contains Summoner Spell information from Riot API
type SummonerSpell struct {
	ID            string             `json:"id"`
	Name          string             `json:"name"`
	Description   string             `json:"description"`
	Tooltip       string             `json:"tooltip"`
	Maxrank       int                `json:"maxrank"`
	Cooldown      []float32          `json:"cooldown"`
	CooldownBurn  string             `json:"cooldownBurn"`
	Cost          []int              `json:"cost"`
	CostBurn      string             `json:"costBurn"`
	Key           string             `json:"key"`
	SummonerLevel int                `json:"summonerLevel"`
	Modes         []string           `json:"modes"`
	CostType      string             `json:"costType"`
	Maxammo       string             `json:"maxammo"`
	Range         []int              `json:"range"`
	RangeBurn     string             `json:"rangeBurn"`
	Image         SummonerSpellImage `json:"image"`
	Resource      string             `json:"resource"`

	Timestamp time.Time `json:"timestamp"`
}

// SummonerSpellsList is used to pass around a list of summoner spells
type SummonerSpellsList map[string]SummonerSpell
