package analyzer

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"

	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/utils"
)

// SingleItemCombiStatistics contains all the statistics for one given item combination described
// by Combi
type SingleItemCombiStatistics struct {
	Combi string // hash describing this combination
	Items []int  // list of all items for this combination

	Picks uint32 // how often was this combination picked
	Wins  uint32 // how often did this combination win
}

// ItemCombiStatistics is a set of SingleItemCombiStatistics which are
// identified by their Combi hash
type ItemCombiStatistics map[string]*SingleItemCombiStatistics // [hash]

// ChampionItemCombiStatistics contains the whole item analysis for a given
// Champion identified by its ID. It contains also the game version for which
// this analysis was performed
type ChampionItemCombiStatistics struct {
	ChampionID int

	GameVersionMajor int
	GameVersionMinor int

	// PerRole is the role in sense of TOP, ..., CARRY, SUPPORT, BOT_UNKNOWN, UNKNOWN
	PerRole           map[string]ItemCombiStatistics // [role]
	PerRoleSampleSize map[string]uint32              // [role]

	Total           ItemCombiStatistics
	TotalSampleSize uint32
}

// ItemAnalyzer is used to perform analysis of best items for a given champion.
// It holds the results and gives back the analzed results if requested.
type ItemAnalyzer struct {
	log *logrus.Entry

	GameVersionMajor int
	GameVersionMinor int

	PerChampion map[int]*ChampionItemCombiStatistics // [ChampionID]
}

// NewItemAnalyzer creates a new champion item analyzer
func NewItemAnalyzer(gameVersionMajor int, gameVersionMinor int) *ItemAnalyzer {
	a := ItemAnalyzer{
		GameVersionMajor: gameVersionMajor,
		GameVersionMinor: gameVersionMinor,
		PerChampion:      make(map[int]*ChampionItemCombiStatistics),

		log: logging.Get(fmt.Sprintf("ItemAnalyzer GameVersion %d.%d", gameVersionMajor, gameVersionMinor)),
	}
	a.log.Trace("New Item Analyzer created")
	return &a
}

func (a *ItemAnalyzer) feedParticipant(p *riotclient.ParticipantDTO) {
	items := []int{
		p.Stats.Item0, p.Stats.Item1, p.Stats.Item2,
		p.Stats.Item3, p.Stats.Item4, p.Stats.Item5,
	}
	for _, item := range items {
		if item == 0 {
			// We only want complete builts
			return
		}
	}
	sort.Ints(items)
	itemCombiHash := utils.HashSortedInt(items)

	championID := p.ChampionID
	role := determineRole(p.Timeline.Lane, p.Timeline.Role)

	// this function will do nothing if it already exists so it should be cheap
	a.addNewRole(championID, role)

	perRole := a.PerChampion[championID].PerRole[role]
	if _, ok := perRole[itemCombiHash]; !ok {
		perRole[itemCombiHash] = &SingleItemCombiStatistics{
			Combi: itemCombiHash,
			Items: items,
		}
	}

	perRole[itemCombiHash].Picks++
	if p.Stats.Win {
		perRole[itemCombiHash].Wins++
	}

	a.PerChampion[championID].PerRoleSampleSize[role]++
	a.PerChampion[championID].PerRole[role] = perRole
}

func determineRole(lane, role string) string {
	switch strings.ToUpper(lane) {
	case "TOP":
		return "TOP"
	case "MID":
		fallthrough
	case "MIDDLE":
		return "MIDDLE"
	case "JUNGLE":
		return "JUNGLE"
	case "BOT":
		fallthrough
	case "BOTTOM":
		switch strings.ToUpper(role) {
		case "DUO_CARRY":
			return "CARRY"
		case "DUO_SUPPORT":
			return "SUPPORT"
		default:
			return "BOTTOM_UNKNOWN"
		}
	default:
		return "UNKNOWN"
	}
}

// FeedMatch is used to feed a new match to add to the analysis to the Analyzer
func (a *ItemAnalyzer) FeedMatch(m *riotclient.MatchDTO) {
	for idx := range m.Participants {
		a.feedParticipant(&m.Participants[idx])
	}
}

func (a *ItemAnalyzer) generateTotal() {
	for _, data := range a.PerChampion {
		data.TotalSampleSize = 0
		total := make(ItemCombiStatistics)

		for _, roleData := range data.PerRole {
			for hash, itemCombi := range roleData {
				if _, ok := total[hash]; !ok {
					total[hash] = &SingleItemCombiStatistics{
						Combi: itemCombi.Combi,
						Items: itemCombi.Items,
					}
				}
				total[hash].Picks = total[hash].Picks + itemCombi.Picks
				total[hash].Wins = total[hash].Wins + itemCombi.Wins

				data.TotalSampleSize += itemCombi.Picks
			}
		}

		data.Total = total
	}
}

// Analyze performs the final analysis and returns the results
func (a *ItemAnalyzer) Analyze() map[int]*ChampionItemCombiStatistics {
	a.generateTotal()
	return a.PerChampion
}

func (a *ItemAnalyzer) addNewChampion(championID int) {
	if _, ok := a.PerChampion[championID]; !ok {
		a.PerChampion[championID] = &ChampionItemCombiStatistics{
			ChampionID:        championID,
			GameVersionMajor:  a.GameVersionMajor,
			GameVersionMinor:  a.GameVersionMinor,
			PerRole:           make(map[string]ItemCombiStatistics),
			PerRoleSampleSize: make(map[string]uint32),
		}
	}
}

func (a *ItemAnalyzer) addNewRole(championID int, role string) {
	a.addNewChampion(championID)
	if _, ok := a.PerChampion[championID].PerRole[role]; !ok {
		a.PerChampion[championID].PerRole[role] = make(ItemCombiStatistics)
	}
	if _, ok := a.PerChampion[championID].PerRoleSampleSize[role]; !ok {
		a.PerChampion[championID].PerRoleSampleSize[role] = 0
	}
}
