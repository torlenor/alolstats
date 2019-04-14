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

func TestRiotClientV4_ActiveGameBySummonerID(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	type args struct {
		summonerID string
	}
	tests := []struct {
		name              string
		fields            fields
		args              args
		want              *riotclient.CurrentGameInfoDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid ActiveGame JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "someEncryptedSummonerID",
			},
			want: &riotclient.CurrentGameInfoDTO{
				GameID:            3875417125,
				GameStartTime:     1545818812438,
				PlatformID:        "EUW1",
				GameMode:          "CLASSIC",
				MapID:             11,
				GameType:          "MATCHED_GAME",
				GameQueueConfigID: 420,
				Observers: riotclient.ObserverDTO{
					EncryptionKey: "ebpxKVFlx4PYsXThI0dgD5TE12tKdpDQ",
				},
				Participants: []riotclient.CurrentGameParticipantDTO{
					{
						ProfileIconID:            3398,
						ChampionID:               53,
						SummonerName:             "g1ngeren",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8400,
							PerkIds: []int64{
								8439,
								8446,
								8429,
								8451,
								8345,
								8347,
								5007,
								5002,
								5001,
							},
							PerkSubStyle: 8300,
						},
						Spell2ID:   14,
						TeamID:     100,
						Spell1ID:   4,
						SummonerID: "N5tdlu92HebFaJ-bYUasQtnt_bAyN4pPkESRQdIUeVDxrJol",
					},
					{
						ProfileIconID:            2072,
						ChampionID:               60,
						SummonerName:             "361ufo",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8100,
							PerkIds: []int64{
								8112,
								8143,
								8120,
								8105,
								9111,
								8014,
								5005,
								5008,
								5002,
							},
							PerkSubStyle: 8000,
						},
						Spell2ID:   11,
						TeamID:     100,
						Spell1ID:   4,
						SummonerID: "S0kbmji0sI-wpBXmt3SW0Tua3YpQQJ2UAzkfAdIxWoWpTduH",
					},
					{
						ProfileIconID:            711,
						ChampionID:               3,
						SummonerName:             "SUP Elramir",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8400,
							PerkIds: []int64{
								8439,
								8401,
								8473,
								8453,
								8226,
								8210,
								5008,
								5008,
								5003,
							},
							PerkSubStyle: 8200,
						},
						Spell2ID:   4,
						TeamID:     100,
						Spell1ID:   12,
						SummonerID: "qvYRubeD30ZYxn3Tq2hlPrOg7nBvu1ckhMkV-hWIAx2f8Ek",
					},
					{
						ProfileIconID:            3505,
						ChampionID:               145,
						SummonerName:             "1907 FB Hades",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8000,
							PerkIds: []int64{
								8005,
								9111,
								9104,
								8014,
								8139,
								8135,
								5005,
								5008,
								5002,
							},
							PerkSubStyle: 8100,
						},
						Spell2ID:   7,
						TeamID:     100,
						Spell1ID:   4,
						SummonerID: "Gy31xT-6ZAITBV3X37yDwYX1WKVQ5iT9EZBCMnukYWku5A5S",
					},
					{
						ProfileIconID:            3303,
						ChampionID:               13,
						SummonerName:             "SPY Vizicsacsi",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8400,
							PerkIds: []int64{
								8439,
								8446,
								8473,
								8451,
								8139,
								8135,
								5008,
								5008,
								5002,
							},
							PerkSubStyle: 8100,
						},
						Spell2ID:   4,
						TeamID:     100,
						Spell1ID:   12,
						SummonerID: "-zdUfLZFUSPv_su82sz1rmzpm-R6A9IDJd6IQEhrdzofhRw",
					},
					{
						ProfileIconID:            3886,
						ChampionID:               69,
						SummonerName:             "Jester 5",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8200,
							PerkIds: []int64{
								8229,
								8226,
								8233,
								8237,
								8139,
								8135,
								5008,
								5008,
								5003,
							},
							PerkSubStyle: 8100,
						},
						Spell2ID:   1,
						TeamID:     200,
						Spell1ID:   4,
						SummonerID: "r-_BJ6SxbQyGctmggO-L09LMDiEe4ACKJ0Ge5HSWAMNO9yQ",
					},
					{
						ProfileIconID:            3624,
						ChampionID:               92,
						SummonerName:             "Bl1ce",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8000,
							PerkIds: []int64{
								8010,
								9111,
								9104,
								8014,
								8275,
								8210,
								5008,
								5008,
								5001,
							},
							PerkSubStyle: 8200,
						},
						Spell2ID:   12,
						TeamID:     200,
						Spell1ID:   4,
						SummonerID: "GFnbAlNLS868cEgOmfeEQaMVjWrugcLHLPMnL1k9bz-66Jo",
					},
					{
						ProfileIconID:            565,
						ChampionID:               89,
						SummonerName:             "lur1keen",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8400,
							PerkIds: []int64{
								8439,
								8463,
								8473,
								8242,
								8306,
								8316,
								5007,
								5002,
								5001,
							},
							PerkSubStyle: 8300,
						},
						Spell2ID:   14,
						TeamID:     200,
						Spell1ID:   4,
						SummonerID: "ToM_r3pNL-89FL4LtVZEGjnLQb4DM-JIx82KaIeqrYaW4OU",
					},
					{
						ProfileIconID:            3797,
						ChampionID:               81,
						SummonerName:             "Gentside xMatty",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8300,
							PerkIds: []int64{
								8359,
								8304,
								8345,
								8347,
								8226,
								8233,
								5008,
								5008,
								5003,
							},
							PerkSubStyle: 8200,
						},
						Spell2ID:   4,
						TeamID:     200,
						Spell1ID:   7,
						SummonerID: "q4wzP8KktDEAEvJaLZrlrEMf1zuNwvg-ljp2EVbpihMCYCw",
					},
					{
						ProfileIconID:            3176,
						ChampionID:               79,
						SummonerName:             "Suirotas",
						GameCustomizationObjects: []riotclient.GameCustomizationObjectDTO{},
						Bot:                      false,
						Perks: riotclient.PerksDTO{
							PerkStyle: 8100,
							PerkIds: []int64{
								8124,
								8143,
								8120,
								8105,
								8233,
								8232,
								5005,
								5008,
								5002,
							},
							PerkSubStyle: 8200,
						},
						Spell2ID:   4,
						TeamID:     200,
						Spell1ID:   11,
						SummonerID: "7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE",
					},
				},
				GameLength: 282,
				BannedChampions: []riotclient.BannedChampionDTO{
					{
						TeamID:     100,
						ChampionID: 39,
						PickTurn:   1,
					},
					{
						TeamID:     100,
						ChampionID: 266,
						PickTurn:   2,
					},
					{
						TeamID:     100,
						ChampionID: 164,
						PickTurn:   3,
					},
					{
						TeamID:     100,
						ChampionID: 121,
						PickTurn:   4,
					},
					{
						TeamID:     100,
						ChampionID: 28,
						PickTurn:   5,
					},
					{
						TeamID:     200,
						ChampionID: 84,
						PickTurn:   6,
					},
					{
						TeamID:     200,
						ChampionID: 236,
						PickTurn:   7,
					},
					{
						TeamID:     200,
						ChampionID: 39,
						PickTurn:   8,
					},
					{
						TeamID:     200,
						ChampionID: 6,
						PickTurn:   9,
					},
					{
						TeamID:     200,
						ChampionID: 12,
						PickTurn:   10,
					},
				},
			},
			wantErr:           false,
			setJSON:           []byte(`{"gameId":3875417125,"gameStartTime":1545818812438,"platformId":"EUW1","gameMode":"CLASSIC","mapId":11,"gameType":"MATCHED_GAME","gameQueueConfigId":420,"observers":{"encryptionKey":"ebpxKVFlx4PYsXThI0dgD5TE12tKdpDQ"},"participants":[{"profileIconId":3398,"ChampionId":53,"summonerName":"g1ngeren","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8446,8429,8451,8345,8347,5007,5002,5001],"perkSubStyle":8300},"spell2Id":14,"TeamID":100,"spell1Id":4,"summonerId":"N5tdlu92HebFaJ-bYUasQtnt_bAyN4pPkESRQdIUeVDxrJol"},{"profileIconId":2072,"ChampionId":60,"summonerName":"361ufo","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8100,"perkIds":[8112,8143,8120,8105,9111,8014,5005,5008,5002],"perkSubStyle":8000},"spell2Id":11,"TeamID":100,"spell1Id":4,"summonerId":"S0kbmji0sI-wpBXmt3SW0Tua3YpQQJ2UAzkfAdIxWoWpTduH"},{"profileIconId":711,"ChampionId":3,"summonerName":"SUP Elramir","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8401,8473,8453,8226,8210,5008,5008,5003],"perkSubStyle":8200},"spell2Id":4,"TeamID":100,"spell1Id":12,"summonerId":"qvYRubeD30ZYxn3Tq2hlPrOg7nBvu1ckhMkV-hWIAx2f8Ek"},{"profileIconId":3505,"ChampionId":145,"summonerName":"1907 FB Hades","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8000,"perkIds":[8005,9111,9104,8014,8139,8135,5005,5008,5002],"perkSubStyle":8100},"spell2Id":7,"TeamID":100,"spell1Id":4,"summonerId":"Gy31xT-6ZAITBV3X37yDwYX1WKVQ5iT9EZBCMnukYWku5A5S"},{"profileIconId":3303,"ChampionId":13,"summonerName":"SPY Vizicsacsi","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8446,8473,8451,8139,8135,5008,5008,5002],"perkSubStyle":8100},"spell2Id":4,"TeamID":100,"spell1Id":12,"summonerId":"-zdUfLZFUSPv_su82sz1rmzpm-R6A9IDJd6IQEhrdzofhRw"},{"profileIconId":3886,"ChampionId":69,"summonerName":"Jester 5","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8200,"perkIds":[8229,8226,8233,8237,8139,8135,5008,5008,5003],"perkSubStyle":8100},"spell2Id":1,"TeamID":200,"spell1Id":4,"summonerId":"r-_BJ6SxbQyGctmggO-L09LMDiEe4ACKJ0Ge5HSWAMNO9yQ"},{"profileIconId":3624,"ChampionId":92,"summonerName":"Bl1ce","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8000,"perkIds":[8010,9111,9104,8014,8275,8210,5008,5008,5001],"perkSubStyle":8200},"spell2Id":12,"TeamID":200,"spell1Id":4,"summonerId":"GFnbAlNLS868cEgOmfeEQaMVjWrugcLHLPMnL1k9bz-66Jo"},{"profileIconId":565,"ChampionId":89,"summonerName":"lur1keen","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8463,8473,8242,8306,8316,5007,5002,5001],"perkSubStyle":8300},"spell2Id":14,"TeamID":200,"spell1Id":4,"summonerId":"ToM_r3pNL-89FL4LtVZEGjnLQb4DM-JIx82KaIeqrYaW4OU"},{"profileIconId":3797,"ChampionId":81,"summonerName":"Gentside xMatty","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8300,"perkIds":[8359,8304,8345,8347,8226,8233,5008,5008,5003],"perkSubStyle":8200},"spell2Id":4,"TeamID":200,"spell1Id":7,"summonerId":"q4wzP8KktDEAEvJaLZrlrEMf1zuNwvg-ljp2EVbpihMCYCw"},{"profileIconId":3176,"ChampionId":79,"summonerName":"Suirotas","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8100,"perkIds":[8124,8143,8120,8105,8233,8232,5005,5008,5002],"perkSubStyle":8200},"spell2Id":4,"TeamID":200,"spell1Id":11,"summonerId":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE"}],"gameLength":282,"bannedChampions":[{"TeamID":100,"ChampionId":39,"PickTurn":1},{"TeamID":100,"ChampionId":266,"PickTurn":2},{"TeamID":100,"ChampionId":164,"PickTurn":3},{"TeamID":100,"ChampionId":121,"PickTurn":4},{"TeamID":100,"ChampionId":28,"PickTurn":5},{"TeamID":200,"ChampionId":84,"PickTurn":6},{"TeamID":200,"ChampionId":236,"PickTurn":7},{"TeamID":200,"ChampionId":39,"PickTurn":8},{"TeamID":200,"ChampionId":6,"PickTurn":9},{"TeamID":200,"ChampionId":12,"PickTurn":10}]}`),
			setError:          nil,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/someEncryptedSummonerID",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 2 - Receive invalid ActiveGame JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "someEncryptedSummonerID2",
			},
			setJSON:           []byte(`{[{{{"gameId":3875417125,"gameStartTime":1545818812438,"platformId":"EUW1","gameMode":"CLASSIC","mapId":11,"gameType":"MATCHED_GAME","gameQueueConfigId":420,"observers":{"encryptionKey":"ebpxKVFlx4PYsXThI0dgD5TE12tKdpDQ"},"participants":[{"profileIconId":3398,"ChampionId":53,"summonerName":"g1ngeren","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8446,8429,8451,8345,8347,5007,5002,5001],"perkSubStyle":8300},"spell2Id":14,"TeamID":100,"spell1Id":4,"summonerId":"N5tdlu92HebFaJ-bYUasQtnt_bAyN4pPkESRQdIUeVDxrJol"},{"profileIconId":2072,"ChampionId":60,"summonerName":"361ufo","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8100,"perkIds":[8112,8143,8120,8105,9111,8014,5005,5008,5002],"perkSubStyle":8000},"spell2Id":11,"TeamID":100,"spell1Id":4,"summonerId":"S0kbmji0sI-wpBXmt3SW0Tua3YpQQJ2UAzkfAdIxWoWpTduH"},{"profileIconId":711,"ChampionId":3,"summonerName":"SUP Elramir","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8401,8473,8453,8226,8210,5008,5008,5003],"perkSubStyle":8200},"spell2Id":4,"TeamID":100,"spell1Id":12,"summonerId":"qvYRubeD30ZYxn3Tq2hlPrOg7nBvu1ckhMkV-hWIAx2f8Ek"},{"profileIconId":3505,"ChampionId":145,"summonerName":"1907 FB Hades","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8000,"perkIds":[8005,9111,9104,8014,8139,8135,5005,5008,5002],"perkSubStyle":8100},"spell2Id":7,"TeamID":100,"spell1Id":4,"summonerId":"Gy31xT-6ZAITBV3X37yDwYX1WKVQ5iT9EZBCMnukYWku5A5S"},{"profileIconId":3303,"ChampionId":13,"summonerName":"SPY Vizicsacsi","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8446,8473,8451,8139,8135,5008,5008,5002],"perkSubStyle":8100},"spell2Id":4,"TeamID":100,"spell1Id":12,"summonerId":"-zdUfLZFUSPv_su82sz1rmzpm-R6A9IDJd6IQEhrdzofhRw"},{"profileIconId":3886,"ChampionId":69,"summonerName":"Jester 5","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8200,"perkIds":[8229,8226,8233,8237,8139,8135,5008,5008,5003],"perkSubStyle":8100},"spell2Id":1,"TeamID":200,"spell1Id":4,"summonerId":"r-_BJ6SxbQyGctmggO-L09LMDiEe4ACKJ0Ge5HSWAMNO9yQ"},{"profileIconId":3624,"ChampionId":92,"summonerName":"Bl1ce","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8000,"perkIds":[8010,9111,9104,8014,8275,8210,5008,5008,5001],"perkSubStyle":8200},"spell2Id":12,"TeamID":200,"spell1Id":4,"summonerId":"GFnbAlNLS868cEgOmfeEQaMVjWrugcLHLPMnL1k9bz-66Jo"},{"profileIconId":565,"ChampionId":89,"summonerName":"lur1keen","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8400,"perkIds":[8439,8463,8473,8242,8306,8316,5007,5002,5001],"perkSubStyle":8300},"spell2Id":14,"TeamID":200,"spell1Id":4,"summonerId":"ToM_r3pNL-89FL4LtVZEGjnLQb4DM-JIx82KaIeqrYaW4OU"},{"profileIconId":3797,"ChampionId":81,"summonerName":"Gentside xMatty","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8300,"perkIds":[8359,8304,8345,8347,8226,8233,5008,5008,5003],"perkSubStyle":8200},"spell2Id":4,"TeamID":200,"spell1Id":7,"summonerId":"q4wzP8KktDEAEvJaLZrlrEMf1zuNwvg-ljp2EVbpihMCYCw"},{"profileIconId":3176,"ChampionId":79,"summonerName":"Suirotas","gameCustomizationObjects":[],"bot":false,"perks":{"perkStyle":8100,"perkIds":[8124,8143,8120,8105,8233,8232,5005,5008,5002],"perkSubStyle":8200},"spell2Id":4,"TeamID":200,"spell1Id":11,"summonerId":"7w1cHKdOPa9XRHe2Lm5x9dBdm1UsbuPRw3FPXm-_O40dykE"}],"gameLength":282,"bannedChampions":[{"TeamID":100,"ChampionId":39,"PickTurn":1},{"TeamID":100,"ChampionId":266,"PickTurn":2},{"TeamID":100,"ChampionId":164,"PickTurn":3},{"TeamID":100,"ChampionId":121,"PickTurn":4},{"TeamID":100,"ChampionId":28,"PickTurn":5},{"TeamID":200,"ChampionId":84,"PickTurn":6},{"TeamID":200,"ChampionId":236,"PickTurn":7},{"TeamID":200,"ChampionId":39,"PickTurn":8},{"TeamID":200,"ChampionId":6,"PickTurn":9},{"TeamID":200,"ChampionId":12,"PickTurn":10}]}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/someEncryptedSummonerID2",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 3 - API Call failure",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			args: args{
				summonerID: "someEncryptedSummonerID3",
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/someEncryptedSummonerID3",
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

			got, err := c.ActiveGameBySummonerID(tt.args.summonerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.ActiveGameBySummonerID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.ActiveGameBySummonerID() = %v, want %v", got, tt.want)
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

func TestRiotClientV4_FeaturedGames(t *testing.T) {
	// Override real API call with our fake one
	apiCall = (*RiotClientV4).mockAPICall

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config config.RiotClient
		log    *logrus.Entry
	}
	tests := []struct {
		name              string
		fields            fields
		want              *riotclient.FeaturedGamesDTO
		wantErr           bool
		setJSON           []byte
		setError          error
		wantAPICallPath   string
		wantAPICallMethod string
		wantAPICallBody   string
	}{
		{
			name: "Test 1 - Receive valid Featured Games JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			want: &riotclient.FeaturedGamesDTO{
				ClientRefreshInterval: 300,
				GameList: []riotclient.FeaturedGameInfoDTO{
					{
						GameID:            3875591604,
						GameStartTime:     1545827415409,
						PlatformID:        "EUW1",
						GameMode:          "ARAM",
						MapID:             12,
						GameType:          "MATCHED_GAME",
						GameQueueConfigID: 450,
						Observers: riotclient.ObserverDTO{
							EncryptionKey: "vqB25+TUzW9yVPxW6vcBt1hmkBY/uMi6",
						},
						Participants: []riotclient.FeaturedGameInfoParticipantDTO{
							{
								ProfileIconID: 661,
								ChampionID:    101,
								SummonerName:  "Mixcharlie",
								Bot:           false,
								Spell2ID:      4,
								TeamID:        100,
								Spell1ID:      3,
							},
							{
								ProfileIconID: 588,
								ChampionID:    222,
								SummonerName:  "3829",
								Bot:           false,
								Spell2ID:      7,
								TeamID:        100,
								Spell1ID:      4,
							},
							{
								ProfileIconID: 1,
								ChampionID:    92,
								SummonerName:  "Desesto",
								Bot:           false,
								Spell2ID:      4,
								TeamID:        100,
								Spell1ID:      32,
							},
							{
								ProfileIconID: 589,
								ChampionID:    18,
								SummonerName:  "BroDa59",
								Bot:           false,
								Spell2ID:      4,
								TeamID:        100,
								Spell1ID:      7,
							},
							{
								ProfileIconID: 3872,
								ChampionID:    126,
								SummonerName:  "Fr34ky95",
								Bot:           false,
								Spell2ID:      32,
								TeamID:        100,
								Spell1ID:      4,
							},
							{
								ProfileIconID: 3854,
								ChampionID:    69,
								SummonerName:  "Cerrolus",
								Bot:           false,
								Spell2ID:      13,
								TeamID:        200,
								Spell1ID:      4,
							},
							{
								ProfileIconID: 3872,
								ChampionID:    42,
								SummonerName:  "CG TroyMcClure",
								Bot:           false,
								Spell2ID:      7,
								TeamID:        200,
								Spell1ID:      4,
							},
							{
								ProfileIconID: 3853,
								ChampionID:    267,
								SummonerName:  "DoomsdayMachine",
								Bot:           false,
								Spell2ID:      13,
								TeamID:        200,
								Spell1ID:      4,
							},
							{
								ProfileIconID: 3505,
								ChampionID:    33,
								SummonerName:  "Naliwaj",
								Bot:           false,
								Spell2ID:      4,
								TeamID:        200,
								Spell1ID:      32,
							},
							{
								ProfileIconID: 3872,
								ChampionID:    45,
								SummonerName:  "letimtimmy",
								Bot:           false,
								Spell2ID:      4,
								TeamID:        200,
								Spell1ID:      7,
							},
						},
						GameLength:      355,
						BannedChampions: []riotclient.BannedChampionDTO{},
					},
				},
			},
			setJSON:           []byte(`{"clientRefreshInterval":300,"gameList":[{"gameId":3875591604,"gameStartTime":1545827415409,"platformId":"EUW1","gameMode":"ARAM","mapId":12,"gameType":"MATCHED_GAME","gameQueueConfigId":450,"observers":{"encryptionKey":"vqB25+TUzW9yVPxW6vcBt1hmkBY/uMi6"},"participants":[{"profileIconId":661,"championId":101,"summonerName":"Mixcharlie","bot":false,"spell2Id":4,"teamId":100,"spell1Id":3},{"profileIconId":588,"championId":222,"summonerName":"3829","bot":false,"spell2Id":7,"teamId":100,"spell1Id":4},{"profileIconId":1,"championId":92,"summonerName":"Desesto","bot":false,"spell2Id":4,"teamId":100,"spell1Id":32},{"profileIconId":589,"championId":18,"summonerName":"BroDa59","bot":false,"spell2Id":4,"teamId":100,"spell1Id":7},{"profileIconId":3872,"championId":126,"summonerName":"Fr34ky95","bot":false,"spell2Id":32,"teamId":100,"spell1Id":4},{"profileIconId":3854,"championId":69,"summonerName":"Cerrolus","bot":false,"spell2Id":13,"teamId":200,"spell1Id":4},{"profileIconId":3872,"championId":42,"summonerName":"CG TroyMcClure","bot":false,"spell2Id":7,"teamId":200,"spell1Id":4},{"profileIconId":3853,"championId":267,"summonerName":"DoomsdayMachine","bot":false,"spell2Id":13,"teamId":200,"spell1Id":4},{"profileIconId":3505,"championId":33,"summonerName":"Naliwaj","bot":false,"spell2Id":4,"teamId":200,"spell1Id":32},{"profileIconId":3872,"championId":45,"summonerName":"letimtimmy","bot":false,"spell2Id":4,"teamId":200,"spell1Id":7}],"gameLength":355,"bannedChampions":[]}]}`),
			setError:          nil,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/featured-games",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},

		{
			name: "Test 2 - Receive invalid Features Games JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			setJSON:           []byte(`{dddd__{}"clientRefreshInterval":300,"gameList":[{"gameId":3875591604,"gameStartTime":1545827415409,"platformId":"EUW1","gameMode":"ARAM","mapId":12,"gameType":"MATCHED_GAME","gameQueueConfigId":450,"observers":{"encryptionKey":"vqB25+TUzW9yVPxW6vcBt1hmkBY/uMi6"},"participants":[{"profileIconId":661,"championId":101,"summonerName":"Mixcharlie","bot":false,"spell2Id":4,"teamId":100,"spell1Id":3},{"profileIconId":588,"championId":222,"summonerName":"3829","bot":false,"spell2Id":7,"teamId":100,"spell1Id":4},{"profileIconId":1,"championId":92,"summonerName":"Desesto","bot":false,"spell2Id":4,"teamId":100,"spell1Id":32},{"profileIconId":589,"championId":18,"summonerName":"BroDa59","bot":false,"spell2Id":4,"teamId":100,"spell1Id":7},{"profileIconId":3872,"championId":126,"summonerName":"Fr34ky95","bot":false,"spell2Id":32,"teamId":100,"spell1Id":4},{"profileIconId":3854,"championId":69,"summonerName":"Cerrolus","bot":false,"spell2Id":13,"teamId":200,"spell1Id":4},{"profileIconId":3872,"championId":42,"summonerName":"CG TroyMcClure","bot":false,"spell2Id":7,"teamId":200,"spell1Id":4},{"profileIconId":3853,"championId":267,"summonerName":"DoomsdayMachine","bot":false,"spell2Id":13,"teamId":200,"spell1Id":4},{"profileIconId":3505,"championId":33,"summonerName":"Naliwaj","bot":false,"spell2Id":4,"teamId":200,"spell1Id":32},{"profileIconId":3872,"championId":45,"summonerName":"letimtimmy","bot":false,"spell2Id":4,"teamId":200,"spell1Id":7}],"gameLength":355,"bannedChampions":[]}]}`),
			setError:          nil,
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/featured-games",
			wantAPICallMethod: "GET",
			wantAPICallBody:   "",
		},
		{
			name: "Test 3 - API Call failure",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
			},
			setJSON:           []byte(""),
			setError:          fmt.Errorf("Some API error"),
			wantErr:           true,
			wantAPICallPath:   "https://euw1.api.riotgames.com/lol/spectator/v4/featured-games",
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

			got, err := c.FeaturedGames()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.FeaturedGames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RiotClientV4.FeaturedGames() = %v, want %v", got, tt.want)
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
