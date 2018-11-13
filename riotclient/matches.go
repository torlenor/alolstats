package riotclient

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Match contains the complete match data (excluding full time line)
type Match struct {
	GameID       int64  `json:"gameId"`
	PlatformID   string `json:"platformId"`
	GameCreation int64  `json:"gameCreation"`
	GameDuration int    `json:"gameDuration"`
	QueueID      int    `json:"queueId"`
	MapID        int    `json:"mapId"`
	SeasonID     int    `json:"seasonId"`
	GameVersion  string `json:"gameVersion"`
	GameMode     string `json:"gameMode"`
	GameType     string `json:"gameType"`
	Teams        []struct {
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
	} `json:"teams"`
	Participants []struct {
		ParticipantID int `json:"participantId"`
		TeamID        int `json:"teamId"`
		ChampionID    int `json:"championId"`
		Spell1ID      int `json:"spell1Id"`
		Spell2ID      int `json:"spell2Id"`
		Masteries     []struct {
			MasteryID int `json:"masteryId"`
			Rank      int `json:"rank"`
		} `json:"masteries"`
		Runes []struct {
			RuneID int `json:"runeId"`
			Rank   int `json:"rank"`
		} `json:"runes"`
		HighestAchievedSeasonTier string `json:"highestAchievedSeasonTier"`
		Stats                     struct {
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
		} `json:"stats"`
		Timeline struct {
			ParticipantID      int `json:"participantId"`
			CreepsPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    int     `json:"20-30"`
			} `json:"creepsPerMinDeltas"`
			XpPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    float64 `json:"20-30"`
			} `json:"xpPerMinDeltas"`
			GoldPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    float64 `json:"20-30"`
			} `json:"goldPerMinDeltas"`
			CsDiffPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    int     `json:"20-30"`
			} `json:"csDiffPerMinDeltas"`
			XpDiffPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    float64 `json:"20-30"`
			} `json:"xpDiffPerMinDeltas"`
			DamageTakenPerMinDeltas struct {
				One020    float64 `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    float64 `json:"20-30"`
			} `json:"damageTakenPerMinDeltas"`
			DamageTakenDiffPerMinDeltas struct {
				One020    int     `json:"10-20"`
				Zero10    float64 `json:"0-10"`
				Three0End float64 `json:"30-end"`
				Two030    float64 `json:"20-30"`
			} `json:"damageTakenDiffPerMinDeltas"`
			Role string `json:"role"`
			Lane string `json:"lane"`
		} `json:"timeline"`
	} `json:"participants"`
	ParticipantIdentities []struct {
		ParticipantID int `json:"participantId"`
		Player        struct {
			PlatformID        string `json:"platformId"`
			AccountID         int    `json:"accountId"`
			SummonerName      string `json:"summonerName"`
			SummonerID        int    `json:"summonerId"`
			CurrentPlatformID string `json:"currentPlatformId"`
			CurrentAccountID  int    `json:"currentAccountId"`
			MatchHistoryURI   string `json:"matchHistoryUri"`
			ProfileIcon       int    `json:"profileIcon"`
		} `json:"player"`
	} `json:"participantIdentities"`
}

// Matches contains a collection of matches
type Matches struct {
	Matches []Match `json:"matches"`
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
	EndIndex   int `json:"endIndex"`
	StartIndex int `json:"startIndex"`
	TotalGames int `json:"totalGames"`
}

// MatchByID gets a match by its ID
func (c *RiotClient) MatchByID(id uint64) (s *Match, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v3/matches/3827449823
	idStr := strconv.FormatUint(id, 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/v3/matches/"+idStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	match := Match{}
	err = json.Unmarshal(data, &match)
	if err != nil {
		return nil, err
	} else if match.GameID == 0 {
		return nil, fmt.Errorf("Match GameID invalid, probably empty data")
	}

	return &match, nil
}

// MatchByAccoundID gets a match by AccountID. Provide a start and end index to fetch matches.
func (c *RiotClient) MatchByAccoundID(id uint64, startIndex uint32, endIndex uint32) (s *MatchList, err error) {
	// Example: https://euw1.api.riotgames.com/lol/match/v3/matchlists/by-account/40722898?beginIndex=100&endIndex=100
	idStr := strconv.FormatUint(id, 10)
	startIndexStr := strconv.FormatUint(uint64(startIndex), 10)
	endIndexStr := strconv.FormatUint(uint64(endIndex), 10)
	data, err := c.apiCall("https://"+c.config.Region+".api.riotgames.com/lol/match/v3/matchlists/by-account/"+idStr+"?beginIndex="+startIndexStr+"?endIndex="+endIndexStr, "GET", "")
	if err != nil {
		return nil, fmt.Errorf("Error in API call: %s", err)
	}

	matchList := MatchList{}
	err = json.Unmarshal(data, &matchList)
	if err != nil {
		return nil, err
	}

	return &matchList, nil
}
