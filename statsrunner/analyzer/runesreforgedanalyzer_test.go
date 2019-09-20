package analyzer

import (
	"reflect"
	"testing"

	"git.abyle.org/hps/alolstats/riotclient"
	"git.abyle.org/hps/alolstats/storage"
	"git.abyle.org/hps/alolstats/utils"
)

func Test_fillRunesReforgedPicks(t *testing.T) {
	type args struct {
		stats riotclient.ParticipantStatsDTO
	}
	tests := []struct {
		name string
		args args
		want storage.RunesReforgedPicks
	}{
		{
			name: "Test 1 - Some runes reforged / perk values",
			args: args{
				stats: riotclient.ParticipantStatsDTO{
					PerkPrimaryStyle: 1,
					Perk0:            2,
					Perk1:            3,
					Perk2:            4,
					Perk3:            5,
					PerkSubStyle:     6,
					Perk4:            7,
					Perk5:            8,
					StatPerk0:        9,
					StatPerk1:        10,
					StatPerk2:        11,
				},
			},
			want: storage.RunesReforgedPicks{
				SlotPrimary: struct {
					storage.Rune

					Rune0 storage.Rune
					Rune1 storage.Rune
					Rune2 storage.Rune
					Rune3 storage.Rune
				}{Rune: storage.Rune{1}, Rune0: storage.Rune{2}, Rune1: storage.Rune{3}, Rune2: storage.Rune{4}, Rune3: storage.Rune{5}},
				SlotSecondary: struct {
					storage.Rune

					Rune0 storage.Rune
					Rune1 storage.Rune
				}{Rune: storage.Rune{6}, Rune0: storage.Rune{7}, Rune1: storage.Rune{8}},
				StatPerks: struct {
					Perk0 storage.Rune
					Perk1 storage.Rune
					Perk2 storage.Rune
				}{Perk0: storage.Rune{9}, Perk1: storage.Rune{10}, Perk2: storage.Rune{11}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fillRunesReforgedPicks(tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("fillRunesReforgedPicks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getRunesReforgedInts(t *testing.T) {
	type args struct {
		stats riotclient.ParticipantStatsDTO
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Test 1 - Some runes reforged / perk values",
			args: args{
				stats: riotclient.ParticipantStatsDTO{
					PerkPrimaryStyle: 1,
					Perk0:            2,
					Perk1:            3,
					Perk2:            4,
					Perk3:            5,
					PerkSubStyle:     6,
					Perk4:            7,
					Perk5:            8,
					StatPerk0:        9,
					StatPerk1:        10,
					StatPerk2:        11,
				},
			},
			want: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
		},
		{
			name: "Test 2 - Some other runes reforged / perk values",
			args: args{
				stats: riotclient.ParticipantStatsDTO{
					PerkPrimaryStyle: 61,
					Perk0:            22,
					Perk1:            73,
					Perk2:            34,
					Perk3:            45,
					PerkSubStyle:     556,
					Perk4:            67,
					Perk5:            78,
					StatPerk0:        99,
					StatPerk1:        10,
					StatPerk2:        111,
				},
			},
			want: []int{61, 22, 73, 34, 45, 556, 67, 78, 99, 10, 111},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getRunesReforgedInts(tt.args.stats); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getRunesReforgedInts() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRunesReforgedAnalyzer_addNewChampion(t *testing.T) {
	type args struct {
		championID int
		lane       string
		role       string
	}

	name := "Test 1 - Create valid champion map when Champion does not exist"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, lane: "top", role: "support"}
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.addNewChampion(arg.championID)
		if _, ok := a.PerChampion[arg.championID]; !ok {
			t.Errorf("championID %d not created", arg.championID)
			return
		}
		if len(a.PerChampion[arg.championID].PerRole) != 0 {
			t.Errorf("expected empty map for lane instead got map of len %d", len(a.PerChampion[arg.championID].PerRole))
			return
		}
	})

	name = "Test 2 - Do not create new champion when map when Champion exists"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, lane: "top", role: "support"}
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.addNewChampion(arg.championID)
		// to make sure the champion is not recreated we fake the id to use as a check
		c := a.PerChampion[arg.championID]
		c.ChampionID = 567
		a.PerChampion[arg.championID] = c

		a.addNewChampion(arg.championID)
		if a.PerChampion[arg.championID].ChampionID != 567 {
			t.Errorf("championID is %d, want %d", a.PerChampion[arg.championID].ChampionID, 567)
			return
		}
		if len(a.PerChampion[arg.championID].PerRole) != 0 {
			t.Errorf("expected empty map for lane instead got map of len %d", len(a.PerChampion[arg.championID].PerRole))
			return
		}
	})
}

func TestRunesReforgedAnalyzer_addNewRoles(t *testing.T) {
	type args struct {
		championID int
		role       string
	}

	name := "Test 1 - Create valid role map when Champion does not exist"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, role: "SUPPORT"}
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.addNewRole(arg.championID, arg.role)
		if _, ok := a.PerChampion[arg.championID]; !ok {
			t.Errorf("championID %d not created", arg.championID)
			return
		}
		if _, ok := a.PerChampion[arg.championID].PerRole[arg.role]; !ok {
			t.Errorf("lane %s not created", arg.role)
			return
		}
		if len(a.PerChampion[arg.championID].PerRole[arg.role]) != 0 {
			t.Errorf("expected empty map for role instead got map of len %d", len(a.PerChampion[arg.championID].PerRole[arg.role]))
			return
		}
	})

	name = "Test 2 - Create valid role map when Champion exists"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, role: "TOP"}
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.addNewChampion(arg.championID)
		// to make sure the champion is not recreated we fake the id to use as a check
		c := a.PerChampion[arg.championID]
		c.ChampionID = 567
		a.PerChampion[arg.championID] = c

		a.addNewRole(arg.championID, arg.role)
		if a.PerChampion[arg.championID].ChampionID != 567 {
			t.Errorf("championID is %d, want %d", a.PerChampion[arg.championID].ChampionID, 567)
			return
		}
		if _, ok := a.PerChampion[arg.championID].PerRole[arg.role]; !ok {
			t.Errorf("lane %s not created", arg.role)
			return
		}
		if len(a.PerChampion[arg.championID].PerRole[arg.role]) != 0 {
			t.Errorf("expected empty map for role instead got map of len %d", len(a.PerChampion[arg.championID].PerRole[arg.role]))
			return
		}
	})
}

func TestRunesReforgedAnalyzer_feedParticipant(t *testing.T) {
	const (
		championID1 = 123
		lane1       = "TOP"
		role1       = "SomeRole"
		detRole1    = "TOP"
	)
	pCombiWin := riotclient.ParticipantDTO{
		ChampionID: championID1,
		Timeline: riotclient.ParticipantTimelineDTO{
			Lane: lane1,
			Role: role1,
		},
		Stats: riotclient.ParticipantStatsDTO{
			PerkPrimaryStyle: 1,
			Perk0:            2,
			Perk1:            3,
			Perk2:            4,
			Perk3:            5,
			PerkSubStyle:     6,
			Perk4:            7,
			Perk5:            8,
			StatPerk0:        9,
			StatPerk1:        10,
			StatPerk2:        11,
			Win:              true,
		},
	}

	pCombiLoss := riotclient.ParticipantDTO{
		ChampionID: championID1,
		Timeline: riotclient.ParticipantTimelineDTO{
			Lane: lane1,
			Role: role1,
		},
		Stats: riotclient.ParticipantStatsDTO{
			PerkPrimaryStyle: 1,
			Perk0:            2,
			Perk1:            3,
			Perk2:            4,
			Perk3:            5,
			PerkSubStyle:     6,
			Perk4:            7,
			Perk5:            8,
			StatPerk0:        9,
			StatPerk1:        10,
			StatPerk2:        11,
			Win:              false,
		},
	}

	name := "Test 1 - Feed one participant"
	t.Run(name, func(t *testing.T) {
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.feedParticipant(&pCombiWin)
		if _, ok := a.PerChampion[championID1]; !ok {
			t.Errorf("championID %d not created", championID1)
			return
		}
		if _, ok := a.PerChampion[championID1].PerRole[detRole1]; !ok {
			t.Errorf("role not created %v", a.PerChampion[championID1].PerRole)
			return
		}
		runesReforgedHash := utils.HashInt(getRunesReforgedInts(pCombiWin.Stats))
		if _, ok := a.PerChampion[championID1].PerRole[detRole1][runesReforgedHash]; !ok {
			t.Errorf("runesReforged combi not created")
			return
		}
		results := a.PerChampion[championID1].PerRole[detRole1][runesReforgedHash]
		if results.Picks != 1 {
			t.Errorf("picks are %d, want %d", results.Picks, 1)
		}
		if results.Wins != 1 {
			t.Errorf("wins are %d, want %d", results.Wins, 1)
		}
		if a.PerChampion[championID1].PerRoleSampleSize[detRole1] != 1 {
			t.Errorf("PerRolePicks are %d, want %d", a.PerChampion[championID1].PerRoleSampleSize[detRole1], 1)
		}
	})

	name = "Test 2 - Feed two participants"
	t.Run(name, func(t *testing.T) {
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		a.feedParticipant(&pCombiWin)
		a.feedParticipant(&pCombiLoss)
		if _, ok := a.PerChampion[championID1]; !ok {
			t.Errorf("championID %d not created", championID1)
			return
		}
		if _, ok := a.PerChampion[championID1].PerRole[detRole1]; !ok {
			t.Errorf("role not created")
			return
		}
		runesReforgedHash := utils.HashInt(getRunesReforgedInts(pCombiWin.Stats))
		if _, ok := a.PerChampion[championID1].PerRole[detRole1][runesReforgedHash]; !ok {
			t.Errorf("runesReforged combi not created")
			return
		}
		results := a.PerChampion[championID1].PerRole[detRole1][runesReforgedHash]
		if results.Picks != 2 {
			t.Errorf("picks are %d, want %d", results.Picks, 2)
		}
		if results.Wins != 1 {
			t.Errorf("wins are %d, want %d", results.Wins, 1)
		}
		if a.PerChampion[championID1].PerRoleSampleSize[detRole1] != 2 {
			t.Errorf("PerRolePicks are %d, want %d", a.PerChampion[championID1].PerRoleSampleSize[detRole1], 2)
		}
	})
}

func TestRunesReforgedAnalyzer_generateTotal(t *testing.T) {

	name := "Test 1 - Generate Total"
	t.Run(name, func(t *testing.T) {
		a := NewRunesReforgedAnalyzer(gameVersionMajor, gameVersionMinor)

		runesReforged1 := riotclient.ParticipantStatsDTO{
			PerkPrimaryStyle: 1,
			Perk0:            2,
			Perk1:            3,
			Perk2:            4,
			Perk3:            5,
			PerkSubStyle:     6,
			Perk4:            7,
			Perk5:            8,
			StatPerk0:        9,
			StatPerk1:        10,
			StatPerk2:        11,
		}
		runesReforged1Filled := fillRunesReforgedPicks(runesReforged1)
		runesReforged1Hash := utils.HashInt(getRunesReforgedInts(runesReforged1))

		runesReforged2 := riotclient.ParticipantStatsDTO{
			PerkPrimaryStyle: 51,
			Perk0:            52,
			Perk1:            53,
			Perk2:            54,
			Perk3:            55,
			PerkSubStyle:     56,
			Perk4:            57,
			Perk5:            58,
			StatPerk0:        59,
			StatPerk1:        510,
			StatPerk2:        511,
		}
		runesReforged2Filled := fillRunesReforgedPicks(runesReforged2)
		runesReforged2Hash := utils.HashInt(getRunesReforgedInts(runesReforged2))

		a.addNewRole(123, "SomeRole1")
		a.PerChampion[123].PerRole["SomeRole1"][runesReforged1Hash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforged1Hash,
			RunesReforged: runesReforged1Filled,
			Picks:         43,
			Wins:          5,
		}
		a.PerChampion[123].PerRole["SomeRole1"][runesReforged2Hash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforged2Hash,
			RunesReforged: runesReforged2Filled,
			Picks:         3,
			Wins:          1,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole1"] = 43 + 3

		a.addNewRole(123, "SomeRole2")
		a.PerChampion[123].PerRole["SomeRole2"][runesReforged1Hash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforged1Hash,
			RunesReforged: runesReforged1Filled,
			Picks:         10,
			Wins:          2,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole2"] = 10

		a.addNewRole(123, "SomeRole3")
		a.PerChampion[123].PerRole["SomeRole3"][runesReforged1Hash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforged1Hash,
			RunesReforged: runesReforged1Filled,
			Picks:         20,
			Wins:          4,
		}
		a.PerChampion[123].PerRole["SomeRole3"][runesReforged2Hash] = &SingleRunesReforgedCombiStatistics{
			Combi:         runesReforged2Hash,
			RunesReforged: runesReforged2Filled,
			Picks:         3,
			Wins:          1,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole3"] = 20 + 3

		a.generateTotal()

		if _, ok := a.PerChampion[123].Total[runesReforged1Hash]; !ok {
			t.Errorf("runesReforged combi %s not in total", runesReforged1Hash)
			return
		}
		if a.PerChampion[123].TotalSampleSize != 43+3+10+20+3 {
			t.Errorf("Total Sample Size not %d but %d", 43+10+20+3, a.PerChampion[123].TotalSampleSize)
			return
		}
		if a.PerChampion[123].PerRoleSampleSize["SomeRole1"] != 43+3 {
			t.Errorf("Sample Size per role/lane not %d but %d", 43+3, a.PerChampion[123].PerRoleSampleSize["SomeRole1"])
			return
		}
		result := a.PerChampion[123].Total[runesReforged1Hash]
		if result.Combi != runesReforged1Hash {
			t.Errorf("Combi field in Total struct not correct, is %s, want %s", result.Combi, runesReforged1Hash)
			return
		}
		if !reflect.DeepEqual(result.RunesReforged, runesReforged1Filled) {
			t.Errorf("RunesReforged slice not as expected, got %v, want %v", result.RunesReforged, runesReforged1Filled)
			return
		}
		if result.Picks != 43+10+20 {
			t.Errorf("Picks are %d, want %d", result.Picks, 43+10+20)
			return
		}
		if result.Wins != 5+2+4 {
			t.Errorf("Wins are %d, want %d", result.Wins, 5+2+4)
		}
	})
}
