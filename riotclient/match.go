package riotclient

// ParticipantStatsDTO contains Stats for a participant of a match
type ParticipantStatsDTO struct {
	ParticipantID                   int  `json:"participantId"`
	Win                             bool `json:"win"`
	Item0                           int  `json:"item0"`
	Item1                           int  `json:"item1"`
	Item2                           int  `json:"item2"`
	Item3                           int  `json:"item3"`
	Item4                           int  `json:"item4"`
	Item5                           int  `json:"item5"`
	Item6                           int  `json:"item6"`
	Kills                           int  `json:"kills"`
	Deaths                          int  `json:"deaths"`
	Assists                         int  `json:"assists"`
	LargestKillingSpree             int  `json:"largestKillingSpree"`
	LargestMultiKill                int  `json:"largestMultiKill"`
	KillingSprees                   int  `json:"killingSprees"`
	LongestTimeSpentLiving          int  `json:"longestTimeSpentLiving"`
	DoubleKills                     int  `json:"doubleKills"`
	TripleKills                     int  `json:"tripleKills"`
	QuadraKills                     int  `json:"quadraKills"`
	PentaKills                      int  `json:"pentaKills"`
	UnrealKills                     int  `json:"unrealKills"`
	TotalDamageDealt                int  `json:"totalDamageDealt"`
	MagicDamageDealt                int  `json:"magicDamageDealt"`
	PhysicalDamageDealt             int  `json:"physicalDamageDealt"`
	TrueDamageDealt                 int  `json:"trueDamageDealt"`
	LargestCriticalStrike           int  `json:"largestCriticalStrike"`
	TotalDamageDealtToChampions     int  `json:"totalDamageDealtToChampions"`
	MagicDamageDealtToChampions     int  `json:"magicDamageDealtToChampions"`
	PhysicalDamageDealtToChampions  int  `json:"physicalDamageDealtToChampions"`
	TrueDamageDealtToChampions      int  `json:"trueDamageDealtToChampions"`
	TotalHeal                       int  `json:"totalHeal"`
	TotalUnitsHealed                int  `json:"totalUnitsHealed"`
	DamageSelfMitigated             int  `json:"damageSelfMitigated"`
	DamageDealtToObjectives         int  `json:"damageDealtToObjectives"`
	DamageDealtToTurrets            int  `json:"damageDealtToTurrets"`
	VisionScore                     int  `json:"visionScore"`
	TimeCCingOthers                 int  `json:"timeCCingOthers"`
	TotalDamageTaken                int  `json:"totalDamageTaken"`
	MagicalDamageTaken              int  `json:"magicalDamageTaken"`
	PhysicalDamageTaken             int  `json:"physicalDamageTaken"`
	TrueDamageTaken                 int  `json:"trueDamageTaken"`
	GoldEarned                      int  `json:"goldEarned"`
	GoldSpent                       int  `json:"goldSpent"`
	TurretKills                     int  `json:"turretKills"`
	InhibitorKills                  int  `json:"inhibitorKills"`
	TotalMinionsKilled              int  `json:"totalMinionsKilled"`
	NeutralMinionsKilled            int  `json:"neutralMinionsKilled"`
	NeutralMinionsKilledTeamJungle  int  `json:"neutralMinionsKilledTeamJungle"`
	NeutralMinionsKilledEnemyJungle int  `json:"neutralMinionsKilledEnemyJungle"`
	TotalTimeCrowdControlDealt      int  `json:"totalTimeCrowdControlDealt"`
	ChampLevel                      int  `json:"champLevel"`
	VisionWardsBoughtInGame         int  `json:"visionWardsBoughtInGame"`
	SightWardsBoughtInGame          int  `json:"sightWardsBoughtInGame"`
	WardsPlaced                     int  `json:"wardsPlaced"`
	WardsKilled                     int  `json:"wardsKilled"`
	FirstBloodKill                  bool `json:"firstBloodKill"`
	FirstBloodAssist                bool `json:"firstBloodAssist"`
	FirstTowerKill                  bool `json:"firstTowerKill"`
	FirstTowerAssist                bool `json:"firstTowerAssist"`
	FirstInhibitorKill              bool `json:"firstInhibitorKill"`
	FirstInhibitorAssist            bool `json:"firstInhibitorAssist"`
	CombatPlayerScore               int  `json:"combatPlayerScore"`
	ObjectivePlayerScore            int  `json:"objectivePlayerScore"`
	TotalPlayerScore                int  `json:"totalPlayerScore"`
	TotalScoreRank                  int  `json:"totalScoreRank"`
	PlayerScore0                    int  `json:"playerScore0"`
	PlayerScore1                    int  `json:"playerScore1"`
	PlayerScore2                    int  `json:"playerScore2"`
	PlayerScore3                    int  `json:"playerScore3"`
	PlayerScore4                    int  `json:"playerScore4"`
	PlayerScore5                    int  `json:"playerScore5"`
	PlayerScore6                    int  `json:"playerScore6"`
	PlayerScore7                    int  `json:"playerScore7"`
	PlayerScore8                    int  `json:"playerScore8"`
	PlayerScore9                    int  `json:"playerScore9"`
}

// RuneDTO contains the runes of a participant of a match
type RuneDTO struct {
	RuneID int `json:"runeId"`
	Rank   int `json:"rank"`
}

// ParticipantTimelineDTO contains participant timeline data for a match
type ParticipantTimelineDTO struct {
	ParticipantID               int                `json:"participantId"`
	CreepsPerMinDeltas          map[string]float64 `json:"creepsPerMinDeltas"`
	XpPerMinDeltas              map[string]float64 `json:"xpPerMinDeltas"`
	GoldPerMinDeltas            map[string]float64 `json:"goldPerMinDeltas"`
	CsDiffPerMinDeltas          map[string]float64 `json:"csDiffPerMinDeltas"`
	XpDiffPerMinDeltas          map[string]float64 `json:"xpDiffPerMinDeltas"`
	DamageTakenPerMinDeltas     map[string]float64 `json:"damageTakenPerMinDeltas"`
	DamageTakenDiffPerMinDeltas map[string]float64 `json:"damageTakenDiffPerMinDeltas"`
	Role                        string             `json:"role"`
	Lane                        string             `json:"lane"`
}

// MasteryDTO list of legacy Mastery information. Not included for matches played with Runes Reforged.
type MasteryDTO struct {
	MasteryID int `json:"masteryId"`
	Rank      int `json:"rank"`
}

// ParticipantDTO contains match data specific to one player
type ParticipantDTO struct {
	ParticipantID             int                    `json:"participantId"`
	TeamID                    int                    `json:"teamId"`
	ChampionID                int                    `json:"championId"`
	Spell1ID                  int                    `json:"spell1Id"`
	Spell2ID                  int                    `json:"spell2Id"`
	Masteries                 []MasteryDTO           `json:"masteries"`
	Runes                     []RuneDTO              `json:"runes"`
	HighestAchievedSeasonTier string                 `json:"highestAchievedSeasonTier"`
	Stats                     ParticipantStatsDTO    `json:"stats"`
	Timeline                  ParticipantTimelineDTO `json:"timeline"`
}

// TeamStatsDTO contains Team information
type TeamStatsDTO struct {
	TeamID               int    `json:"teamId"`
	Win                  string `json:"win"`
	FirstBlood           bool   `json:"firstBlood"`
	FirstTower           bool   `json:"firstTower"`
	FirstInhibitor       bool   `json:"firstInhibitor"`
	FirstBaron           bool   `json:"firstBaron"`
	FirstDragon          bool   `json:"firstDragon"`
	FirstRiftHerald      bool   `json:"firstRiftHerald"`
	TowerKills           int    `json:"towerKills"`
	InhibitorKills       int    `json:"inhibitorKills"`
	BaronKills           int    `json:"baronKills"`
	DragonKills          int    `json:"dragonKills"`
	VilemawKills         int    `json:"vilemawKills"`
	RiftHeraldKills      int    `json:"riftHeraldKills"`
	DominionVictoryScore int    `json:"dominionVictoryScore"`
	Bans                 []struct {
		ChampionID int `json:"championId"`
		PickTurn   int `json:"pickTurn"`
	} `json:"bans"`
}

// PlayerDTO contains the Player Information
type PlayerDTO struct {
	PlatformID        string `json:"platformId"`
	AccountID         string `json:"accountId"`
	SummonerName      string `json:"summonerName"`
	SummonerID        string `json:"summonerId"`
	CurrentPlatformID string `json:"currentPlatformId"`
	CurrentAccountID  string `json:"currentAccountId"`
	MatchHistoryURI   string `json:"matchHistoryUri"`
	ProfileIcon       int    `json:"profileIcon"`
}

// ParticipantIdentityDTO contains Participant identity information
type ParticipantIdentityDTO struct {
	ParticipantID int       `json:"participantId"`
	Player        PlayerDTO `json:"player"`
}

// MatchDTO contains the complete match data (excluding full time line)
type MatchDTO struct {
	GameID                int64                    `json:"gameId"`
	PlatformID            string                   `json:"platformId"`
	GameCreation          int64                    `json:"gameCreation"`
	GameDuration          int                      `json:"gameDuration"`
	QueueID               int                      `json:"queueId"`
	MapID                 int                      `json:"mapId"`
	SeasonID              int                      `json:"seasonId"`
	GameVersion           string                   `json:"gameVersion"`
	GameMode              string                   `json:"gameMode"`
	GameType              string                   `json:"gameType"`
	Teams                 []TeamStatsDTO           `json:"teams"`
	Participants          []ParticipantDTO         `json:"participants" bson:"participants"`
	ParticipantIdentities []ParticipantIdentityDTO `json:"participantIdentities"`
}

// Matches contains a collection of MatchDTOs
type Matches struct {
	Matches []MatchDTO `json:"matches"`
}

// MatchList contains a list of matches which have been requested via the API
type MatchList struct {
	Matches []struct {
		Lane       string `json:"lane"`
		GameID     int64  `json:"gameId"`
		Champion   int    `json:"champion"`
		PlatformID string `json:"platformId"`
		Timestamp  int64  `json:"timestamp"`
		Queue      int    `json:"queue"`
		Role       string `json:"role"`
		Season     int    `json:"season"`
	} `json:"matches"`
	StartIndex int `json:"startIndex"`
	EndIndex   int `json:"endIndex"`
	TotalGames int `json:"totalGames"`
}

// MatchReferenceDTO models the Riot API response for one Matches entry in a MatchlistDto
type MatchReferenceDTO struct {
	Lane       string `json:"lane"`
	GameID     int64  `json:"gameId"`
	Champion   int    `json:"champion"`
	PlatformID string `json:"platformId"`
	Timestamp  int64  `json:"timestamp"`
	Queue      int    `json:"queue"`
	Role       string `json:"role"`
	Season     int    `json:"season"`
}

// MatchlistDTO models the Riot API response for Matchlist endpoints
type MatchlistDTO struct {
	Matches    []MatchReferenceDTO `json:"matches"`
	TotalGames int                 `json:"totalGames"`
	StartIndex int                 `json:"startIndex"`
	EndIndex   int                 `json:"endIndex"`
}
