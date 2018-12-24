package riotclient

// CurrentGameInfoDTO contains live game information of a currently running match
type CurrentGameInfoDTO struct {
	GameID            int64  `json:"gameId"`
	GameStartTime     int64  `json:"gameStartTime"`
	PlatformID        string `json:"platformId"`
	GameMode          string `json:"gameMode"`
	MapID             int64  `json:"mapId"`
	GameType          string `json:"gameType"`
	GameQueueConfigID int64  `json:"gameQueueConfigId"`
	BannedChampions   []struct {
		TeamID     int64 `json:"teamId"`
		ChampionID int64 `json:"championId"`
		PickTurn   int   `json:"pickTurn"`
	} `json:"bannedChampions"`
	Observers struct {
		EncryptionKey string `json:"encryptionKey"`
	} `json:"observers"`
	Participants []struct {
		ProfileIconID            int64  `json:"profileIconId"`
		ChampionID               int64  `json:"championId"`
		SummonerName             string `json:"summonerName"`
		GameCustomizationObjects []struct {
			Category string `json:"category"`
			Content  string `json:"content"`
		} `json:"gameCustomizationObjects"`
		Bot   bool `json:"bot"`
		Perks struct {
			PerkStyle    int64   `json:"perkStyle"`
			PerkIds      []int64 `json:"perkIds"`
			PerkSubStyle int64   `json:"perkSubStyle"`
		} `json:"perks"`
		Spell2ID   int64  `json:"spell2Id"`
		TeamID     int64  `json:"teamId"`
		Spell1ID   int64  `json:"spell1Id"`
		SummonerID string `json:"summonerId"`
	} `json:"participants"`
	GameLength int64 `json:"gameLength"`
}
