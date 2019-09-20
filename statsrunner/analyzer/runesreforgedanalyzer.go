package analyzer

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/storage"
	"git.abyle.org/hps/alolstats/utils"
)

// SingleRunesReforgedCombiStatistics contains all the statistics for one given item combination
// described by Combi,
type SingleRunesReforgedCombiStatistics struct {
	Combi         string // hash describing this combination
	RunesReforged storage.RunesReforgedPicks

	Picks uint32 // how often was this combination picked
	Wins  uint32 // how often did this combination win
}

// RunesReforgedCombiStatistics is a set of SingleRunesReforgedCombiStatistics which are
// identified by their Combi hash,
type RunesReforgedCombiStatistics map[string]*SingleRunesReforgedCombiStatistics // [hash]

// ChampionRunesReforgedCombiStatistics contains the whole item analysis for a given
// Champion identified by its ID. It contains also the game version for which
// this analysis was performed.
type ChampionRunesReforgedCombiStatistics struct {
	ChampionID int

	GameVersionMajor int
	GameVersionMinor int

	// PerRole is the role in sense of TOP, ..., CARRY, SUPPORT, BOT_UNKNOWN, UNKNOWN
	PerRole           map[string]RunesReforgedCombiStatistics // [role]
	PerRoleSampleSize map[string]uint32                       // [role]

	Total           RunesReforgedCombiStatistics
	TotalSampleSize uint32
}

// RunesReforgedAnalyzer is used to perform analysis of best items for a given champion.
// It holds the results and gives back the analzed results if requested.
type RunesReforgedAnalyzer struct {
	log *logrus.Entry

	GameVersionMajor int
	GameVersionMinor int

	PerChampion map[int]*ChampionRunesReforgedCombiStatistics // [ChampionID]
}

// NewRunesReforgedAnalyzer creates a new champion item analyzer
func NewRunesReforgedAnalyzer(gameVersionMajor int, gameVersionMinor int) *RunesReforgedAnalyzer {
	a := RunesReforgedAnalyzer{
		GameVersionMajor: gameVersionMajor,
		GameVersionMinor: gameVersionMinor,
		PerChampion:      make(map[int]*ChampionRunesReforgedCombiStatistics),

		log: logging.Get(fmt.Sprintf("RunesReforgedAnalyzer GameVersion %d.%d", gameVersionMajor, gameVersionMinor)),
	}
	a.log.Trace("New RunesReforged Analyzer created")
	return &a
}

func fillRunesReforgedPicks(stats riotclient.ParticipantStatsDTO) storage.RunesReforgedPicks {
	runesReforgedPicks := storage.RunesReforgedPicks{}
	runesReforgedPicks.SlotPrimary.ID = stats.PerkPrimaryStyle
	runesReforgedPicks.SlotPrimary.Rune0.ID = stats.Perk0
	runesReforgedPicks.SlotPrimary.Rune1.ID = stats.Perk1
	runesReforgedPicks.SlotPrimary.Rune2.ID = stats.Perk2
	runesReforgedPicks.SlotPrimary.Rune3.ID = stats.Perk3

	runesReforgedPicks.SlotSecondary.ID = stats.PerkSubStyle
	runesReforgedPicks.SlotSecondary.Rune0.ID = stats.Perk4
	runesReforgedPicks.SlotSecondary.Rune1.ID = stats.Perk5

	runesReforgedPicks.StatPerks.Perk0.ID = stats.StatPerk0
	runesReforgedPicks.StatPerks.Perk1.ID = stats.StatPerk1
	runesReforgedPicks.StatPerks.Perk2.ID = stats.StatPerk2

	return runesReforgedPicks
}

func (a *RunesReforgedAnalyzer) feedParticipant(p *riotclient.ParticipantDTO) {
	runesReforgedForHash := []int{
		p.Stats.PerkPrimaryStyle,
		p.Stats.Perk0,
		p.Stats.Perk1,
		p.Stats.Perk2,
		p.Stats.Perk3,
		p.Stats.PerkSubStyle,
		p.Stats.Perk4,
		p.Stats.Perk5,
		p.Stats.StatPerk0,
		p.Stats.StatPerk1,
		p.Stats.StatPerk2,
	}
	runesReforgedHash := utils.HashInt(runesReforgedForHash)

	runesReforgedPicks := fillRunesReforgedPicks(p.Stats)

	championID := p.ChampionID
	role := determineRole(p.Timeline.Lane, p.Timeline.Role)

	// this function will do nothing if it already exists so it should be cheap
	a.addNewRole(championID, role)

	perRole := a.PerChampion[championID].PerRole[role]
	if _, ok := perRole[runesReforgedHash]; !ok {
		perRole[runesReforgedHash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforgedHash,
			RunesReforged: runesReforgedPicks,
		}
	}

	perRole[runesReforgedHash].Picks++
	if p.Stats.Win {
		perRole[runesReforgedHash].Wins++
	}

	a.PerChampion[championID].PerRoleSampleSize[role]++
	a.PerChampion[championID].PerRole[role] = perRole
}

// FeedMatch is used to feed a new match to add to the analysis to the Analyzer
func (a *RunesReforgedAnalyzer) FeedMatch(m *riotclient.MatchDTO) {
	for idx := range m.Participants {
		a.feedParticipant(&m.Participants[idx])
	}
}

func (a *RunesReforgedAnalyzer) generateTotal() {
	for _, data := range a.PerChampion {
		data.TotalSampleSize = 0
		total := make(RunesReforgedCombiStatistics)

		for _, roleData := range data.PerRole {
			for hash, runesReforgedCombi := range roleData {
				if _, ok := total[hash]; !ok {
					total[hash] = &SingleRunesReforgedCombiStatistics{
						Combi:         runesReforgedCombi.Combi,
						RunesReforged: runesReforgedCombi.RunesReforged,
					}
				}
				total[hash].Picks = total[hash].Picks + runesReforgedCombi.Picks
				total[hash].Wins = total[hash].Wins + runesReforgedCombi.Wins

				data.TotalSampleSize += runesReforgedCombi.Picks
			}
		}

		data.Total = total
	}
}

// Analyze performs the final analysis and returns the results
func (a *RunesReforgedAnalyzer) Analyze() map[int]*ChampionRunesReforgedCombiStatistics {
	a.generateTotal()
	return a.PerChampion
}

func (a *RunesReforgedAnalyzer) addNewChampion(championID int) {
	if _, ok := a.PerChampion[championID]; !ok {
		a.PerChampion[championID] = &ChampionRunesReforgedCombiStatistics{
			ChampionID:        championID,
			GameVersionMajor:  a.GameVersionMajor,
			GameVersionMinor:  a.GameVersionMinor,
			PerRole:           make(map[string]RunesReforgedCombiStatistics),
			PerRoleSampleSize: make(map[string]uint32),
		}
	}
}

func (a *RunesReforgedAnalyzer) addNewRole(championID int, role string) {
	a.addNewChampion(championID)
	if _, ok := a.PerChampion[championID].PerRole[role]; !ok {
		a.PerChampion[championID].PerRole[role] = make(RunesReforgedCombiStatistics)
	}
	if _, ok := a.PerChampion[championID].PerRoleSampleSize[role]; !ok {
		a.PerChampion[championID].PerRoleSampleSize[role] = 0
	}
}
