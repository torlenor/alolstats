package riotclientv4

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/riotclient"
)

func TestRiotClientV4_MatchByID(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.0000"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		id uint64
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.MatchDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Match JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				id: 3872223341,
			},
			wantS: &riotclient.MatchDTO{
				SeasonID: 11,
				QueueID:  400,
				GameID:   3872223341,
				ParticipantIdentities: []riotclient.ParticipantIdentityDTO{
					{
						Player: riotclient.PlayerDTO{
							CurrentPlatformID: "EUW1",
							SummonerName:      "Torlenor",
							MatchHistoryURI:   "/v1/stats/player_history/EUW1/40722898",
							PlatformID:        "EUW1",
							CurrentAccountID:  "1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g",
							ProfileIcon:       987,
							SummonerID:        "BUGNsWBwI7ecsTJ_RjE9_RoZwbiV076vEE2UuLqkVIWyHm0",
							AccountID:         "1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g",
						},
						ParticipantID: 1,
					},
					{
						Player: riotclient.PlayerDTO{
							CurrentPlatformID: "EUW1",
							SummonerName:      "andomore",
							MatchHistoryURI:   "/v1/stats/player_history/EUW1/42708129",
							PlatformID:        "EUW1",
							CurrentAccountID:  "tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q",
							ProfileIcon:       28,
							SummonerID:        "xnL-1jYUSkYQlAOZjmBXMRKEAwWXqZTP0TcrAPKqVKRnkMQ",
							AccountID:         "tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q",
						},
						ParticipantID: 2,
					},
				},
				GameVersion: "8.24.255.8524",
				PlatformID:  "EUW1",
				GameMode:    "CLASSIC",
				MapID:       11,
				GameType:    "MATCHED_GAME",
				Teams: []riotclient.TeamStatsDTO{
					{
						FirstDragon: true,
						Bans: []riotclient.TeamStatsBansDTO{
							{
								PickTurn:   1,
								ChampionID: 238,
							},
							{
								PickTurn:   2,
								ChampionID: 63,
							},
						},
						FirstInhibitor:       false,
						Win:                  "Fail",
						FirstRiftHerald:      true,
						FirstBaron:           false,
						BaronKills:           0,
						RiftHeraldKills:      1,
						FirstBlood:           false,
						TeamID:               100,
						FirstTower:           false,
						VilemawKills:         0,
						InhibitorKills:       0,
						TowerKills:           1,
						DominionVictoryScore: 0,
						DragonKills:          2,
					},
					{
						FirstDragon: false,
						Bans: []riotclient.TeamStatsBansDTO{
							{
								PickTurn:   6,
								ChampionID: 17,
							},
							{
								PickTurn:   7,
								ChampionID: 157,
							},
						},
						FirstInhibitor:       true,
						Win:                  "Win",
						FirstRiftHerald:      false,
						FirstBaron:           false,
						BaronKills:           0,
						RiftHeraldKills:      0,
						FirstBlood:           true,
						TeamID:               200,
						FirstTower:           true,
						VilemawKills:         0,
						InhibitorKills:       2,
						TowerKills:           11,
						DominionVictoryScore: 0,
						DragonKills:          2,
					},
				},
				Participants: []riotclient.ParticipantDTO{
					{
						Stats: riotclient.ParticipantStatsDTO{
							NeutralMinionsKilledTeamJungle:  0,
							VisionScore:                     37,
							MagicDamageDealtToChampions:     21537,
							LargestMultiKill:                1,
							TotalTimeCrowdControlDealt:      193,
							LongestTimeSpentLiving:          1087,
							Perk1Var1:                       250,
							Perk1Var3:                       0,
							Perk1Var2:                       1565,
							TripleKills:                     0,
							Perk5:                           8347,
							Perk4:                           8313,
							PlayerScore9:                    0,
							PlayerScore8:                    0,
							Kills:                           6,
							PlayerScore1:                    0,
							PlayerScore0:                    0,
							PlayerScore3:                    0,
							PlayerScore2:                    0,
							PlayerScore5:                    0,
							PlayerScore4:                    0,
							PlayerScore7:                    0,
							PlayerScore6:                    0,
							Perk5Var1:                       0,
							Perk5Var3:                       0,
							Perk5Var2:                       0,
							TotalScoreRank:                  0,
							NeutralMinionsKilled:            0,
							StatPerk1:                       5002,
							StatPerk0:                       5008,
							DamageDealtToTurrets:            526,
							PhysicalDamageDealtToChampions:  1743,
							DamageDealtToObjectives:         1261,
							Perk2Var2:                       0,
							Perk2Var3:                       0,
							TotalUnitsHealed:                1,
							Perk2Var1:                       20,
							Perk4Var1:                       0,
							TotalDamageTaken:                13340,
							Perk4Var3:                       0,
							WardsKilled:                     1,
							LargestCriticalStrike:           0,
							LargestKillingSpree:             2,
							QuadraKills:                     0,
							MagicDamageDealt:                68578,
							FirstBloodAssist:                false,
							Item2:                           2065,
							Item3:                           3117,
							Item0:                           3157,
							Item1:                           3092,
							Item6:                           3364,
							Item4:                           3065,
							Item5:                           0,
							Perk1:                           8226,
							Perk0:                           8214,
							Perk3:                           8237,
							Perk2:                           8210,
							Perk3Var3:                       0,
							Perk3Var2:                       0,
							Perk3Var1:                       791,
							DamageSelfMitigated:             8820,
							MagicalDamageTaken:              4671,
							Perk0Var2:                       214,
							FirstInhibitorKill:              false,
							TrueDamageTaken:                 299,
							Assists:                         7,
							Perk4Var2:                       0,
							GoldSpent:                       10575,
							TrueDamageDealt:                 3546,
							ParticipantID:                   1,
							PhysicalDamageDealt:             10161,
							SightWardsBoughtInGame:          0,
							TotalDamageDealtToChampions:     24464,
							PhysicalDamageTaken:             8369,
							TotalPlayerScore:                0,
							Win:                             false,
							ObjectivePlayerScore:            0,
							TotalDamageDealt:                82286,
							NeutralMinionsKilledEnemyJungle: 0,
							Deaths:                          4,
							WardsPlaced:                     15,
							PerkPrimaryStyle:                8200,
							PerkSubStyle:                    8300,
							TurretKills:                     0,
							FirstBloodKill:                  false,
							TrueDamageDealtToChampions:      1183,
							GoldEarned:                      10904,
							KillingSprees:                   2,
							UnrealKills:                     0,
							FirstTowerAssist:                false,
							FirstTowerKill:                  false,
							ChampLevel:                      15,
							DoubleKills:                     0,
							InhibitorKills:                  0,
							FirstInhibitorAssist:            false,
							Perk0Var1:                       2296,
							CombatPlayerScore:               0,
							Perk0Var3:                       0,
							VisionWardsBoughtInGame:         1,
							PentaKills:                      0,
							TotalHeal:                       3082,
							TotalMinionsKilled:              62,
							TimeCCingOthers:                 110,
							StatPerk2:                       5003,
						},
						Spell1ID:                  4,
						ParticipantID:             1,
						HighestAchievedSeasonTier: "UNRANKED",
						Spell2ID:                  14,
						TeamID:                    100,
						Timeline: riotclient.ParticipantTimelineDTO{
							Lane:          "BOTTOM",
							ParticipantID: 1,
							CsDiffPerMinDeltas: map[string]float64{
								"20-30": 1.65,
								"0-10":  -0.15000000000000013,
								"10-20": 0.7,
							},
							GoldPerMinDeltas: map[string]float64{
								"20-30": 427.5,
								"0-10":  157.8,
								"10-20": 387.8,
							},
							XpDiffPerMinDeltas: map[string]float64{
								"20-30": 42.80000000000001,
								"0-10":  24.750000000000014,
								"10-20": -0.44999999999998863,
							},
							CreepsPerMinDeltas: map[string]float64{
								"20-30": 3.0999999999999996,
								"0-10":  0.7,
								"10-20": 1.7000000000000002,
							},
							XpPerMinDeltas: map[string]float64{
								"20-30": 547,
								"0-10":  312.6,
								"10-20": 402.3,
							},
							Role: "DUO_SUPPORT",
							DamageTakenDiffPerMinDeltas: map[string]float64{
								"20-30": -0.6500000000000341,
								"0-10":  -190.14999999999998,
								"10-20": -195,
							},
							DamageTakenPerMinDeltas: map[string]float64{
								"20-30": 427.5,
								"0-10":  137.6,
								"10-20": 457,
							},
						},
						ChampionID: 25,
					},
				},
				GameDuration: 1923,
				GameCreation: 1545504274736,
			},
			setJSON:           []byte(`{"seasonId":11,"queueId":400,"gameId":3872223341,"participantIdentities":[{"player":{"currentPlatformId":"EUW1","summonerName":"Torlenor","matchHistoryUri":"/v1/stats/player_history/EUW1/40722898","platformId":"EUW1","currentAccountId":"1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g","profileIcon":987,"summonerId":"BUGNsWBwI7ecsTJ_RjE9_RoZwbiV076vEE2UuLqkVIWyHm0","accountId":"1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g"},"participantId":1},{"player":{"currentPlatformId":"EUW1","summonerName":"andomore","matchHistoryUri":"/v1/stats/player_history/EUW1/42708129","platformId":"EUW1","currentAccountId":"tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q","profileIcon":28,"summonerId":"xnL-1jYUSkYQlAOZjmBXMRKEAwWXqZTP0TcrAPKqVKRnkMQ","accountId":"tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q"},"participantId":2}],"gameVersion":"8.24.255.8524","platformId":"EUW1","gameMode":"CLASSIC","mapId":11,"gameType":"MATCHED_GAME","teams":[{"firstDragon":true,"bans":[{"pickTurn":1,"championId":238},{"pickTurn":2,"championId":63}],"firstInhibitor":false,"win":"Fail","firstRiftHerald":true,"firstBaron":false,"baronKills":0,"riftHeraldKills":1,"firstBlood":false,"teamId":100,"firstTower":false,"vilemawKills":0,"inhibitorKills":0,"towerKills":1,"dominionVictoryScore":0,"dragonKills":2},{"firstDragon":false,"bans":[{"pickTurn":6,"championId":17},{"pickTurn":7,"championId":157}],"firstInhibitor":true,"win":"Win","firstRiftHerald":false,"firstBaron":false,"baronKills":0,"riftHeraldKills":0,"firstBlood":true,"teamId":200,"firstTower":true,"vilemawKills":0,"inhibitorKills":2,"towerKills":11,"dominionVictoryScore":0,"dragonKills":2}],"participants":[{"stats":{"neutralMinionsKilledTeamJungle":0,"visionScore":37,"magicDamageDealtToChampions":21537,"largestMultiKill":1,"totalTimeCrowdControlDealt":193,"longestTimeSpentLiving":1087,"perk1Var1":250,"perk1Var3":0,"perk1Var2":1565,"tripleKills":0,"perk5":8347,"perk4":8313,"playerScore9":0,"playerScore8":0,"kills":6,"playerScore1":0,"playerScore0":0,"playerScore3":0,"playerScore2":0,"playerScore5":0,"playerScore4":0,"playerScore7":0,"playerScore6":0,"perk5Var1":0,"perk5Var3":0,"perk5Var2":0,"totalScoreRank":0,"neutralMinionsKilled":0,"statPerk1":5002,"statPerk0":5008,"damageDealtToTurrets":526,"physicalDamageDealtToChampions":1743,"damageDealtToObjectives":1261,"perk2Var2":0,"perk2Var3":0,"totalUnitsHealed":1,"perk2Var1":20,"perk4Var1":0,"totalDamageTaken":13340,"perk4Var3":0,"wardsKilled":1,"largestCriticalStrike":0,"largestKillingSpree":2,"quadraKills":0,"magicDamageDealt":68578,"firstBloodAssist":false,"item2":2065,"item3":3117,"item0":3157,"item1":3092,"item6":3364,"item4":3065,"item5":0,"perk1":8226,"perk0":8214,"perk3":8237,"perk2":8210,"perk3Var3":0,"perk3Var2":0,"perk3Var1":791,"damageSelfMitigated":8820,"magicalDamageTaken":4671,"perk0Var2":214,"firstInhibitorKill":false,"trueDamageTaken":299,"assists":7,"perk4Var2":0,"goldSpent":10575,"trueDamageDealt":3546,"participantId":1,"physicalDamageDealt":10161,"sightWardsBoughtInGame":0,"totalDamageDealtToChampions":24464,"physicalDamageTaken":8369,"totalPlayerScore":0,"win":false,"objectivePlayerScore":0,"totalDamageDealt":82286,"neutralMinionsKilledEnemyJungle":0,"deaths":4,"wardsPlaced":15,"perkPrimaryStyle":8200,"perkSubStyle":8300,"turretKills":0,"firstBloodKill":false,"trueDamageDealtToChampions":1183,"goldEarned":10904,"killingSprees":2,"unrealKills":0,"firstTowerAssist":false,"firstTowerKill":false,"champLevel":15,"doubleKills":0,"inhibitorKills":0,"firstInhibitorAssist":false,"perk0Var1":2296,"combatPlayerScore":0,"perk0Var3":0,"visionWardsBoughtInGame":1,"pentaKills":0,"totalHeal":3082,"totalMinionsKilled":62,"timeCCingOthers":110,"statPerk2":5003},"spell1Id":4,"participantId":1,"highestAchievedSeasonTier":"UNRANKED","spell2Id":14,"teamId":100,"timeline":{"lane":"BOTTOM","participantId":1,"csDiffPerMinDeltas":{"20-30":1.65,"0-10":-0.15000000000000013,"10-20":0.7},"goldPerMinDeltas":{"20-30":427.5,"0-10":157.8,"10-20":387.8},"xpDiffPerMinDeltas":{"20-30":42.80000000000001,"0-10":24.750000000000014,"10-20":-0.44999999999998863},"creepsPerMinDeltas":{"20-30":3.0999999999999996,"0-10":0.7,"10-20":1.7000000000000002},"xpPerMinDeltas":{"20-30":547,"0-10":312.6,"10-20":402.3},"role":"DUO_SUPPORT","damageTakenDiffPerMinDeltas":{"20-30":-0.6500000000000341,"0-10":-190.14999999999998,"10-20":-195},"damageTakenPerMinDeltas":{"20-30":427.5,"0-10":137.6,"10-20":457}},"championId":25}],"gameDuration":1923,"gameCreation":1545504274736}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matches/3872223341",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid Match JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				id: 3872223341,
			},
			wantS:             nil,
			setJSON:           []byte(`{{{{"seasonId":11,"queueId":400,"gameId":3872223341,"participantIdentities":[{"player":{"currentPlatformId":"EUW1","summonerName":"Torlenor","matchHistoryUri":"/v1/stats/player_history/EUW1/40722898","platformId":"EUW1","currentAccountId":"1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g","profileIcon":987,"summonerId":"BUGNsWBwI7ecsTJ_RjE9_RoZwbiV076vEE2UuLqkVIWyHm0","accountId":"1boL9yr2g5kZbPExCP4I6ngN2NIQxe-gi6FWIC8_Di7D4g"},"participantId":1},{"player":{"currentPlatformId":"EUW1","summonerName":"andomore","matchHistoryUri":"/v1/stats/player_history/EUW1/42708129","platformId":"EUW1","currentAccountId":"tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q","profileIcon":28,"summonerId":"xnL-1jYUSkYQlAOZjmBXMRKEAwWXqZTP0TcrAPKqVKRnkMQ","accountId":"tI7kxDKVBP3iCQLzRGeLmo0z13UzNVCN52kXIJWepys42Q"},"participantId":2}],"gameVersion":"8.24.255.8524","platformId":"EUW1","gameMode":"CLASSIC","mapId":11,"gameType":"MATCHED_GAME","teams":[{"firstDragon":true,"bans":[{"pickTurn":1,"championId":238},{"pickTurn":2,"championId":63}],"firstInhibitor":false,"win":"Fail","firstRiftHerald":true,"firstBaron":false,"baronKills":0,"riftHeraldKills":1,"firstBlood":false,"teamId":100,"firstTower":false,"vilemawKills":0,"inhibitorKills":0,"towerKills":1,"dominionVictoryScore":0,"dragonKills":2},{"firstDragon":false,"bans":[{"pickTurn":6,"championId":17},{"pickTurn":7,"championId":157}],"firstInhibitor":true,"win":"Win","firstRiftHerald":false,"firstBaron":false,"baronKills":0,"riftHeraldKills":0,"firstBlood":true,"teamId":200,"firstTower":true,"vilemawKills":0,"inhibitorKills":2,"towerKills":11,"dominionVictoryScore":0,"dragonKills":2}],"participants":[{"stats":{"neutralMinionsKilledTeamJungle":0,"visionScore":37,"magicDamageDealtToChampions":21537,"largestMultiKill":1,"totalTimeCrowdControlDealt":193,"longestTimeSpentLiving":1087,"perk1Var1":250,"perk1Var3":0,"perk1Var2":1565,"tripleKills":0,"perk5":8347,"perk4":8313,"playerScore9":0,"playerScore8":0,"kills":6,"playerScore1":0,"playerScore0":0,"playerScore3":0,"playerScore2":0,"playerScore5":0,"playerScore4":0,"playerScore7":0,"playerScore6":0,"perk5Var1":0,"perk5Var3":0,"perk5Var2":0,"totalScoreRank":0,"neutralMinionsKilled":0,"statPerk1":5002,"statPerk0":5008,"damageDealtToTurrets":526,"physicalDamageDealtToChampions":1743,"damageDealtToObjectives":1261,"perk2Var2":0,"perk2Var3":0,"totalUnitsHealed":1,"perk2Var1":20,"perk4Var1":0,"totalDamageTaken":13340,"perk4Var3":0,"wardsKilled":1,"largestCriticalStrike":0,"largestKillingSpree":2,"quadraKills":0,"magicDamageDealt":68578,"firstBloodAssist":false,"item2":2065,"item3":3117,"item0":3157,"item1":3092,"item6":3364,"item4":3065,"item5":0,"perk1":8226,"perk0":8214,"perk3":8237,"perk2":8210,"perk3Var3":0,"perk3Var2":0,"perk3Var1":791,"damageSelfMitigated":8820,"magicalDamageTaken":4671,"perk0Var2":214,"firstInhibitorKill":false,"trueDamageTaken":299,"assists":7,"perk4Var2":0,"goldSpent":10575,"trueDamageDealt":3546,"participantId":1,"physicalDamageDealt":10161,"sightWardsBoughtInGame":0,"totalDamageDealtToChampions":24464,"physicalDamageTaken":8369,"totalPlayerScore":0,"win":false,"objectivePlayerScore":0,"totalDamageDealt":82286,"neutralMinionsKilledEnemyJungle":0,"deaths":4,"wardsPlaced":15,"perkPrimaryStyle":8200,"perkSubStyle":8300,"turretKills":0,"firstBloodKill":false,"trueDamageDealtToChampions":1183,"goldEarned":10904,"killingSprees":2,"unrealKills":0,"firstTowerAssist":false,"firstTowerKill":false,"champLevel":15,"doubleKills":0,"inhibitorKills":0,"firstInhibitorAssist":false,"perk0Var1":2296,"combatPlayerScore":0,"perk0Var3":0,"visionWardsBoughtInGame":1,"pentaKills":0,"totalHeal":3082,"totalMinionsKilled":62,"timeCCingOthers":110,"statPerk2":5003},"spell1Id":4,"participantId":1,"highestAchievedSeasonTier":"UNRANKED","spell2Id":14,"teamId":100,"timeline":{"lane":"BOTTOM","participantId":1,"csDiffPerMinDeltas":{"20-30":1.65,"0-10":-0.15000000000000013,"10-20":0.7},"goldPerMinDeltas":{"20-30":427.5,"0-10":157.8,"10-20":387.8},"xpDiffPerMinDeltas":{"20-30":42.80000000000001,"0-10":24.750000000000014,"10-20":-0.44999999999998863},"creepsPerMinDeltas":{"20-30":3.0999999999999996,"0-10":0.7,"10-20":1.7000000000000002},"xpPerMinDeltas":{"20-30":547,"0-10":312.6,"10-20":402.3},"role":"DUO_SUPPORT","damageTakenDiffPerMinDeltas":{"20-30":-0.6500000000000341,"0-10":-190.14999999999998,"10-20":-195},"damageTakenPerMinDeltas":{"20-30":427.5,"0-10":137.6,"10-20":457}},"championId":25}],"gameDuration":1923,"gameCreation":1545504274736}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matches/3872223341",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 3 - API call error",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				id: 3872223341,
			},
			wantS:             nil,
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matches/3872223341",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config: tt.fields.config,
				log:    tt.fields.log,
			}

			apiCallReturnJSON = tt.setJSON
			apiCallReturnErr = tt.setError

			gotS, err := c.MatchByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.MatchByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.MatchByID() = %v, want %v", gotS, tt.wantS)
			}

			if lastAPICallPath != tt.wantAPICallPath {
				t.Errorf("lastAPICallPath = %v, want %v", lastAPICallPath, tt.wantAPICallPath)
			}
			if lastAPICallBody != tt.wantAPICallBody {
				t.Errorf("lastAPICallBody = %v, want %v", lastAPICallBody, tt.wantAPICallBody)
			}
			if lastAPICallMethod != tt.wantAPICallMethod {
				t.Errorf("lastAPICallMethod = %v, want %v", lastAPICallMethod, tt.wantAPICallMethod)
			}
		})
	}
}

func TestRiotClientV4_MatchesByAccountID(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.0000"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		accountID string
		args      map[string]string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantS             *riotclient.MatchlistDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid MatchList JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				args:      make(map[string]string),
			},
			wantS: &riotclient.MatchlistDTO{
				Matches: []riotclient.MatchReferenceDTO{
					{
						Lane:       "JUNGLE",
						GameID:     3875655954,
						Champion:   120,
						PlatformID: "EUW1",
						Timestamp:  1545831637628,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
					{
						Lane:       "JUNGLE",
						GameID:     3875554798,
						Champion:   79,
						PlatformID: "EUW1",
						Timestamp:  1545828172417,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
					{
						Lane:       "JUNGLE",
						GameID:     3875459851,
						Champion:   120,
						PlatformID: "EUW1",
						Timestamp:  1545825706777,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
				},
				EndIndex:   100,
				StartIndex: 0,
				TotalGames: 3,
			},
			setJSON:           []byte(`{"matches":[{"lane":"JUNGLE","gameId":3875655954,"champion":120,"platformId":"EUW1","timestamp":1545831637628,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875554798,"champion":79,"platformId":"EUW1","timestamp":1545828172417,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875459851,"champion":120,"platformId":"EUW1","timestamp":1545825706777,"queue":420,"role":"NONE","season":11}],"endIndex":100,"startIndex":0,"totalGames":3}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid MatchList JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				args:      make(map[string]string),
			},
			wantS:             nil,
			setJSON:           []byte(`{"matches":[{{{{"lane":"JUNGLE","gameId":3875655954,"champion":120,"platformId":"EUW1","timestamp":1545831637628,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875554798,"champion":79,"platformId":"EUW1","timestamp":1545828172417,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875459851,"champion":120,"platformId":"EUW1","timestamp":1545825706777,"queue":420,"role":"NONE","season":11}],"endIndex":100,"startIndex":0,"totalGames":3}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 3 - API error",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				args:      make(map[string]string),
			},
			wantS:             nil,
			setJSON:           []byte(``),
			setError:          fmt.Errorf("Some error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 4 - Set various query args",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				accountID: "C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw",
				args: map[string]string{
					"beginIndex": "0",
					"endIndex":   "100",
					"beginTime":  "0",
					"champion":   "123",
					"queue":      "100",
				},
			},
			wantS: &riotclient.MatchlistDTO{
				Matches: []riotclient.MatchReferenceDTO{
					{
						Lane:       "JUNGLE",
						GameID:     3875655954,
						Champion:   120,
						PlatformID: "EUW1",
						Timestamp:  1545831637628,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
					{
						Lane:       "JUNGLE",
						GameID:     3875554798,
						Champion:   79,
						PlatformID: "EUW1",
						Timestamp:  1545828172417,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
					{
						Lane:       "JUNGLE",
						GameID:     3875459851,
						Champion:   120,
						PlatformID: "EUW1",
						Timestamp:  1545825706777,
						Queue:      420,
						Role:       "NONE",
						Season:     11,
					},
				},
				EndIndex:   100,
				StartIndex: 0,
				TotalGames: 3,
			},
			setJSON:           []byte(`{"matches":[{"lane":"JUNGLE","gameId":3875655954,"champion":120,"platformId":"EUW1","timestamp":1545831637628,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875554798,"champion":79,"platformId":"EUW1","timestamp":1545828172417,"queue":420,"role":"NONE","season":11},{"lane":"JUNGLE","gameId":3875459851,"champion":120,"platformId":"EUW1","timestamp":1545825706777,"queue":420,"role":"NONE","season":11}],"endIndex":100,"startIndex":0,"totalGames":3}`),
			setError:          nil,
			wantErr:           false,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/match/v4/matchlists/by-account/C9VDk9h0oZtvFNWWeQVaU2G_Kq6YWYR2pcKbhmd4TgSMvw?beginIndex=0&beginTime=0&champion=123&endIndex=100&queue=100",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config: tt.fields.config,
				log:    tt.fields.log,
			}

			apiCallReturnJSON = tt.setJSON
			apiCallReturnErr = tt.setError

			gotS, err := c.MatchesByAccountID(tt.args.accountID, tt.args.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.MatchesByAccountID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.MatchesByAccountID() = %v, want %v", gotS, tt.wantS)
			}

			if lastAPICallPath != tt.wantAPICallPath {
				t.Errorf("lastAPICallPath = %v, want %v", lastAPICallPath, tt.wantAPICallPath)
			}
			if lastAPICallBody != tt.wantAPICallBody {
				t.Errorf("lastAPICallBody = %v, want %v", lastAPICallBody, tt.wantAPICallBody)
			}
			if lastAPICallMethod != tt.wantAPICallMethod {
				t.Errorf("lastAPICallMethod = %v, want %v", lastAPICallMethod, tt.wantAPICallMethod)
			}
		})
	}
}

func TestRiotClientV4_MatchTimeLineByID(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.0000"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		matchID uint64
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		wantT             *riotclient.MatchTimelineDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config: tt.fields.config,
				log:    tt.fields.log,
			}

			apiCallReturnJSON = tt.setJSON
			apiCallReturnErr = tt.setError

			gotT, err := c.MatchTimeLineByID(tt.args.matchID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.MatchTimeLineByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotT, tt.wantT) {
				t.Errorf("RiotClientV4.MatchTimeLineByID() = %v, want %v", gotT, tt.wantT)
			}

			if lastAPICallPath != tt.wantAPICallPath {
				t.Errorf("lastAPICallPath = %v, want %v", lastAPICallPath, tt.wantAPICallPath)
			}
			if lastAPICallBody != tt.wantAPICallBody {
				t.Errorf("lastAPICallBody = %v, want %v", lastAPICallBody, tt.wantAPICallBody)
			}
			if lastAPICallMethod != tt.wantAPICallMethod {
				t.Errorf("lastAPICallMethod = %v, want %v", lastAPICallMethod, tt.wantAPICallMethod)
			}
		})
	}
}
