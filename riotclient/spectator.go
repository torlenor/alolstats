package riotclient

// BannedChampionDTO contain the banned champions in a currently running match
type BannedChampionDTO struct {
	TeamID     int64 `json:"teamId"`
	ChampionID int64 `json:"championId"`
	PickTurn   int   `json:"pickTurn"`
}

// ObserverDTO contain the encryption key in a currently running match
type ObserverDTO struct {
	EncryptionKey string `json:"encryptionKey"`
}

// GameCustomizationObjectDTO contains the custom options set for a running match
type GameCustomizationObjectDTO struct {
	Category string `json:"category"`
	Content  string `json:"content"`
}

// PerksDTO contains Perks/Runes Reforged Information for a running match
type PerksDTO struct {
	PerkStyle    int64   `json:"perkStyle"`
	PerkIds      []int64 `json:"perkIds"`
	PerkSubStyle int64   `json:"perkSubStyle"`
}

// CurrentGameParticipantDTO contain information of the participants in the currently running match
type CurrentGameParticipantDTO struct {
	ProfileIconID            int64                        `json:"profileIconId"`
	ChampionID               int64                        `json:"championId"`
	SummonerName             string                       `json:"summonerName"`
	GameCustomizationObjects []GameCustomizationObjectDTO `json:"gameCustomizationObjects"`
	Bot                      bool                         `json:"bot"`
	Perks                    PerksDTO                     `json:"perks"`
	Spell2ID                 int64                        `json:"spell2Id"`
	TeamID                   int64                        `json:"teamId"`
	Spell1ID                 int64                        `json:"spell1Id"`
	SummonerID               string                       `json:"summonerId"`
}

// CurrentGameInfoDTO contains live game information of a currently running match
type CurrentGameInfoDTO struct {
	GameID            int64                       `json:"gameId"`
	GameStartTime     int64                       `json:"gameStartTime"`
	PlatformID        string                      `json:"platformId"`
	GameMode          string                      `json:"gameMode"`
	MapID             int64                       `json:"mapId"`
	GameType          string                      `json:"gameType"`
	GameQueueConfigID int64                       `json:"gameQueueConfigId"`
	BannedChampions   []BannedChampionDTO         `json:"bannedChampions"`
	Observers         ObserverDTO                 `json:"observers"`
	Participants      []CurrentGameParticipantDTO `json:"participants"`
	GameLength        int64                       `json:"gameLength"`
}

// FeaturedGameInfoParticipantDTO contains participant information in a featured game
type FeaturedGameInfoParticipantDTO struct {
	ProfileIconID int    `json:"profileIconId"`
	ChampionID    int    `json:"championId"`
	SummonerName  string `json:"summonerName"`
	Bot           bool   `json:"bot"`
	Spell2ID      int    `json:"spell2Id"`
	TeamID        int    `json:"teamId"`
	Spell1ID      int    `json:"spell1Id"`
}

// FeaturedGameInfoDTO contains informations on one Featured Game
type FeaturedGameInfoDTO struct {
	GameID            int64                            `json:"gameId"`
	GameStartTime     int64                            `json:"gameStartTime"`
	PlatformID        string                           `json:"platformId"`
	GameMode          string                           `json:"gameMode"`
	MapID             int                              `json:"mapId"`
	GameType          string                           `json:"gameType"`
	GameQueueConfigID int                              `json:"gameQueueConfigId"`
	Observers         ObserverDTO                      `json:"observers"`
	Participants      []FeaturedGameInfoParticipantDTO `json:"participants"`
	GameLength        int                              `json:"gameLength"`
	BannedChampions   []BannedChampionDTO              `json:"bannedChampions"`
}

// FeaturedGamesDTO models the Riot API response for Features Games
type FeaturedGamesDTO struct {
	ClientRefreshInterval int                   `json:"clientRefreshInterval"`
	GameList              []FeaturedGameInfoDTO `json:"gameList"`
}
