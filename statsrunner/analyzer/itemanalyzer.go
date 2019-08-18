package analyzer

import (
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

	PerRole map[string]map[string]ItemCombiStatistics // [lane][role]

	Total ItemCombiStatistics
}

// ItemAnalyzer is used to perform analysis of best items for a given champion.
// It holds the results and gives back the analzed results if requested.
type ItemAnalyzer struct {
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
	}
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
	itemCombiHash := utils.HashSortedInt(items)

	role := p.Timeline.Role
	lane := p.Timeline.Lane
	championID := p.ChampionID

	// this function will do nothing if it already exists so it should be cheap
	a.addNewRole(championID, lane, role)

	perRole := a.PerChampion[championID].PerRole[lane][role]
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

	a.PerChampion[championID].PerRole[lane][role] = perRole
}

// FeedMatch is used to feed a new match to add to the analysis to the Analyzer
func (a *ItemAnalyzer) FeedMatch(m *riotclient.MatchDTO) {
	for idx := range m.Participants {
		a.feedParticipant(&m.Participants[idx])
	}
}

func (a *ItemAnalyzer) generateTotal() {

	for _, data := range a.PerChampion {
		total := make(ItemCombiStatistics)

		for _, laneData := range data.PerRole {
			for _, roleData := range laneData {
				for hash, itemCombi := range roleData {
					if _, ok := total[hash]; !ok {
						total[hash] = &SingleItemCombiStatistics{
							Combi: itemCombi.Combi,
							Items: itemCombi.Items,
						}
					}
					total[hash].Picks = total[hash].Picks + itemCombi.Picks
					total[hash].Wins = total[hash].Wins + itemCombi.Wins
				}
			}
		}

		data.Total = total
	}
}

// Analyze performs the final analyzis and returns the results
func (a *ItemAnalyzer) Analyze() map[int]*ChampionItemCombiStatistics {
	a.generateTotal()
	return a.PerChampion
}

func (a *ItemAnalyzer) addNewChampion(championID int) {
	if _, ok := a.PerChampion[championID]; !ok {
		a.PerChampion[championID] = &ChampionItemCombiStatistics{
			ChampionID:       championID,
			GameVersionMajor: a.GameVersionMajor,
			GameVersionMinor: a.GameVersionMinor,
			PerRole:          make(map[string]map[string]ItemCombiStatistics),
		}
	}
}

func (a *ItemAnalyzer) addNewLane(championID int, lane string) {
	a.addNewChampion(championID)
	if _, ok := a.PerChampion[championID].PerRole[lane]; !ok {
		a.PerChampion[championID].PerRole[lane] = make(map[string]ItemCombiStatistics)
	}
}

func (a *ItemAnalyzer) addNewRole(championID int, lane string, role string) {
	a.addNewLane(championID, lane)
	laneEntry := a.PerChampion[championID].PerRole[lane]
	if _, ok := laneEntry[role]; !ok {
		laneEntry[role] = ItemCombiStatistics{}
	}

}
