package analyzer

import (
	"testing"

	"git.abyle.org/hps/alolstats/riotclient"
)

const (
	gameVersionMajor = 9
	gameVersionMinor = 12
)

func TestItemAnalyzer_addNewChampion(t *testing.T) {
	type args struct {
		championID int
		lane       string
		role       string
	}

	name := "Test 1 - Create valid champion map when Champion does not exist"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, lane: "top", role: "support"}
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

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
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

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

func TestItemAnalyzer_addNewLane(t *testing.T) {
	type args struct {
		championID int
		role       string
	}

	name := "Test 1 - Create valid role map when Champion does not exist"
	t.Run(name, func(t *testing.T) {
		arg := args{championID: 123, role: "SUPPORT"}
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

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
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

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

func TestItemAnalyzer_feedParticipant(t *testing.T) {
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
			Item0: 10,
			Item1: 11,
			Item2: 12,
			Item3: 13,
			Item4: 14,
			Item5: 15,
			Win:   true,
		},
	}

	pCombiLoss := riotclient.ParticipantDTO{
		ChampionID: championID1,
		Timeline: riotclient.ParticipantTimelineDTO{
			Lane: lane1,
			Role: role1,
		},
		Stats: riotclient.ParticipantStatsDTO{
			Item0: 10,
			Item1: 11,
			Item2: 12,
			Item3: 13,
			Item4: 14,
			Item5: 15,
			Win:   false,
		},
	}

	pComb2NotEnoughItems := riotclient.ParticipantDTO{
		ChampionID: championID1,
		Timeline: riotclient.ParticipantTimelineDTO{
			Lane: lane1,
			Role: role1,
		},
		Stats: riotclient.ParticipantStatsDTO{
			Item0: 10,
			Item1: 11,
			Win:   false,
		},
	}

	name := "Test 1 - Feed one participant"
	t.Run(name, func(t *testing.T) {
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

		a.feedParticipant(&pCombiWin)
		if _, ok := a.PerChampion[championID1]; !ok {
			t.Errorf("championID %d not created", championID1)
			return
		}
		if _, ok := a.PerChampion[championID1].PerRole[detRole1]; !ok {
			t.Errorf("role not created %v", a.PerChampion[championID1].PerRole)
			return
		}
		if _, ok := a.PerChampion[championID1].PerRole[detRole1]["10_11_12_13_14_15"]; !ok {
			t.Errorf("item combi not created")
			return
		}
		results := a.PerChampion[championID1].PerRole[detRole1]["10_11_12_13_14_15"]
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
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

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
		if _, ok := a.PerChampion[championID1].PerRole[detRole1]["10_11_12_13_14_15"]; !ok {
			t.Errorf("item combi not created")
			return
		}
		results := a.PerChampion[championID1].PerRole[detRole1]["10_11_12_13_14_15"]
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

	name = "Test 3 - Feed participant with not enough items"
	t.Run(name, func(t *testing.T) {
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

		a.feedParticipant(&pComb2NotEnoughItems)
		if _, ok := a.PerChampion[championID1]; ok {
			t.Errorf("championID %d created", championID1)
			return
		}
	})
}

func testEq(a, b []int) bool {
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func TestItemAnalyzer_generateTotal(t *testing.T) {

	name := "Test 1 - Generate Total"
	t.Run(name, func(t *testing.T) {
		a := NewItemAnalyzer(gameVersionMajor, gameVersionMinor)

		items := []int{11, 12, 13, 14, 15, 16}

		a.addNewRole(123, "SomeRole1")
		a.PerChampion[123].PerRole["SomeRole1"]["11_12_13_14_15_16"] = &SingleItemCombiStatistics{
			Combi: "11_12_13_14_15_16",
			Items: items,
			Picks: 43,
			Wins:  5,
		}
		a.PerChampion[123].PerRole["SomeRole1"]["11_12_13_14_15_79"] = &SingleItemCombiStatistics{
			Combi: "11_12_13_14_15_79",
			Items: []int{11, 12, 13, 14, 15, 79},
			Picks: 3,
			Wins:  1,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole1"] = 43 + 3

		a.addNewRole(123, "SomeRole2")
		a.PerChampion[123].PerRole["SomeRole2"]["11_12_13_14_15_16"] = &SingleItemCombiStatistics{
			Combi: "11_12_13_14_15_16",
			Items: items,
			Picks: 10,
			Wins:  2,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole2"] = 10

		a.addNewRole(123, "SomeRole3")
		a.PerChampion[123].PerRole["SomeRole3"]["11_12_13_14_15_16"] = &SingleItemCombiStatistics{
			Combi: "11_12_13_14_15_16",
			Items: items,
			Picks: 20,
			Wins:  4,
		}
		a.PerChampion[123].PerRole["SomeRole3"]["11_12_13_14_15_79"] = &SingleItemCombiStatistics{
			Combi: "11_12_13_14_15_79",
			Items: []int{11, 12, 13, 14, 15, 79},
			Picks: 3,
			Wins:  1,
		}
		a.PerChampion[123].PerRoleSampleSize["SomeRole3"] = 20 + 3

		a.generateTotal()

		if _, ok := a.PerChampion[123].Total["11_12_13_14_15_16"]; !ok {
			t.Errorf("item combi %s not in total", "11_12_13_14_15_16")
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
		result := a.PerChampion[123].Total["11_12_13_14_15_16"]
		if result.Combi != "11_12_13_14_15_16" {
			t.Errorf("Combi field in Total struct not correct, is %s, want %s", result.Combi, "11_12_13_14_15_16")
			return
		}
		if !testEq(result.Items, items) {
			t.Errorf("Items slice not as expected, got %v, want %v", result.Items, items)
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

func Test_determineRole(t *testing.T) {
	type args struct {
		lane string
		role string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test 1 - Determine TOP",
			args: args{
				lane: "TOP",
				role: "something",
			},
			want: "TOP",
		},
		{
			name: "Test 2 - Determine TOP - mixed case",
			args: args{
				lane: "ToP",
				role: "something",
			},
			want: "TOP",
		},
		{
			name: "Test 3 - Determine MIDDLE from MID",
			args: args{
				lane: "MID",
				role: "something",
			},
			want: "MIDDLE",
		},
		{
			name: "Test 4 - Determine MIDDLE from MIDDLE",
			args: args{
				lane: "MIDDLE",
				role: "whatever",
			},
			want: "MIDDLE",
		},
		{
			name: "Test 5 - Determine JUNGLE",
			args: args{
				lane: "JUNGLE",
				role: "whatever",
			},
			want: "JUNGLE",
		},
		{
			name: "Test 6 - Determine CARRY from BOT - DUO_CARRY",
			args: args{
				lane: "BOT",
				role: "DUO_CARRY",
			},
			want: "CARRY",
		},
		{
			name: "Test 7 - Determine SUPPORT from BOTTOM - DUO_carry",
			args: args{
				lane: "BOTTOM",
				role: "DUO_carry",
			},
			want: "CARRY",
		},
		{
			name: "Test 8 - Determine BOTTOM_UNKNOWN",
			args: args{
				lane: "BOTTOM",
				role: "whatever",
			},
			want: "BOTTOM_UNKNOWN",
		},
		{
			name: "Test 9 - Determine UNKNOWN - everything else must be unknown",
			args: args{
				lane: "something",
				role: "whatever",
			},
			want: "UNKNOWN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := determineRole(tt.args.lane, tt.args.role); got != tt.want {
				t.Errorf("determineRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
