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

func TestRiotClientV4_SummonerSpells(t *testing.T) {
	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	type fields struct {
		config  config.RiotClient
		log     *logrus.Entry
		ddragon *MockRiotClientDD
	}
	tests := []struct {
		name    string
		fields  fields
		wantS   *riotclient.SummonerSpellsList
		wantErr bool
	}{
		{
			name: "Test 1 - Receive valid Summoner Spells JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					summonerSpellsJSON: []byte(`{"type":"summoner","version":"9.7.1","data":{"SummonerFlash":{"id":"SummonerFlash","name":"Flash","description":"Teleports your champion a short distance toward your cursor's location.","tooltip":"Teleports your champion a short distance toward your cursor's location.","maxrank":1,"cooldown":[300],"cooldownBurn":"300","cost":[0],"costBurn":"0","datavalues":{},"effect":[null,[400],[0],[0],[0],[0],[0],[0],[0],[0],[0]],"effectBurn":[null,"400","0","0","0","0","0","0","0","0","0"],"vars":[],"key":"4","summonerLevel":7,"modes":["CLASSIC","ODIN","TUTORIAL","ARAM","ASCENSION","FIRSTBLOOD","ASSASSINATE","URF","ARSR","DOOMBOTSTEEMO","STARGUARDIAN","PROJECT","SNOWURF","ONEFORALL","GAMEMODEX","PRACTICETOOL"],"costType":"No Cost","maxammo":"-1","range":[425],"rangeBurn":"425","image":{"full":"SummonerFlash.png","sprite":"spell0.png","group":"spell","x":288,"y":0,"w":48,"h":48},"resource":"No Cost"},"SummonerHeal":{"id":"SummonerHeal","name":"Heal","description":"Restores 90-345 Health (depending on champion level) and grants 30% Movement Speed for 1 second to you and target allied champion. This healing is halved for units recently affected by Summoner Heal.","tooltip":"Restores {{ f1 }} Health and grants 30% Movement Speed for 1 second to your champion and target allied champion. This healing is halved for units recently affected by Summoner Heal.<br /><br /><span class=\"colorFFFF00\">If this spell cannot find a target, it will cast on the most wounded allied champion in range.</span>","maxrank":1,"cooldown":[240],"cooldownBurn":"240","cost":[0],"costBurn":"0","datavalues":{},"effect":[null,[0.3],[75],[15],[0.5],[826],[0.5],[0],[0],[0],[0]],"effectBurn":[null,"0.3","75","15","0.5","826","0.5","0","0","0","0"],"vars":[{"link":"@player.level","coeff":[90,105,120,135,150,165,180,195,210,225,240,255,270,285,300,315,330,345],"key":"f1"}],"key":"7","summonerLevel":1,"modes":["CLASSIC","ODIN","TUTORIAL","ARAM","ASCENSION","FIRSTBLOOD","ASSASSINATE","URF","ARSR","DOOMBOTSTEEMO","STARGUARDIAN","PROJECT","ONEFORALL","TUTORIAL_MODULE_1","TUTORIAL_MODULE_2","GAMEMODEX","PRACTICETOOL"],"costType":"No Cost","maxammo":"-1","range":[850],"rangeBurn":"850","image":{"full":"SummonerHeal.png","sprite":"spell0.png","group":"spell","x":384,"y":0,"w":48,"h":48},"resource":"No Cost"}}}`),
				},
			},
			wantS: &riotclient.SummonerSpellsList{
				"SummonerFlash": {
					Timestamp: now(),

					ID:          "SummonerFlash",
					Name:        "Flash",
					Description: "Teleports your champion a short distance toward your cursor's location.",
					Tooltip:     "Teleports your champion a short distance toward your cursor's location.",
					Maxrank:     1,
					Cooldown: []float32{
						300,
					},
					CooldownBurn: "300",
					Cost: []int{
						0,
					},
					CostBurn:      "0",
					Key:           "4",
					SummonerLevel: 7,
					Modes: []string{
						"CLASSIC",
						"ODIN",
						"TUTORIAL",
						"ARAM",
						"ASCENSION",
						"FIRSTBLOOD",
						"ASSASSINATE",
						"URF",
						"ARSR",
						"DOOMBOTSTEEMO",
						"STARGUARDIAN",
						"PROJECT",
						"SNOWURF",
						"ONEFORALL",
						"GAMEMODEX",
						"PRACTICETOOL",
					},
					CostType: "No Cost",
					Maxammo:  "-1",
					Range: []int{
						425,
					},
					RangeBurn: "425",
					Image: riotclient.SummonerSpellImage{
						Full:   "SummonerFlash.png",
						Sprite: "spell0.png",
						Group:  "spell",
						X:      288,
						Y:      0,
						W:      48,
						H:      48,
					},
					Resource: "No Cost",
				},
				"SummonerHeal": {
					Timestamp: now(),

					ID:          "SummonerHeal",
					Name:        "Heal",
					Description: "Restores 90-345 Health (depending on champion level) and grants 30% Movement Speed for 1 second to you and target allied champion. This healing is halved for units recently affected by Summoner Heal.",
					Tooltip:     "Restores {{ f1 }} Health and grants 30% Movement Speed for 1 second to your champion and target allied champion. This healing is halved for units recently affected by Summoner Heal.<br /><br /><span class=\"colorFFFF00\">If this spell cannot find a target, it will cast on the most wounded allied champion in range.</span>",
					Maxrank:     1,
					Cooldown: []float32{
						240,
					},
					CooldownBurn: "240",
					Cost: []int{
						0,
					},
					CostBurn:      "0",
					Key:           "7",
					SummonerLevel: 1,
					Modes: []string{
						"CLASSIC",
						"ODIN",
						"TUTORIAL",
						"ARAM",
						"ASCENSION",
						"FIRSTBLOOD",
						"ASSASSINATE",
						"URF",
						"ARSR",
						"DOOMBOTSTEEMO",
						"STARGUARDIAN",
						"PROJECT",
						"ONEFORALL",
						"TUTORIAL_MODULE_1",
						"TUTORIAL_MODULE_2",
						"GAMEMODEX",
						"PRACTICETOOL",
					},
					CostType: "No Cost",
					Maxammo:  "-1",
					Range: []int{
						850,
					},
					RangeBurn: "850",
					Image: riotclient.SummonerSpellImage{
						Full:   "SummonerHeal.png",
						Sprite: "spell0.png",
						Group:  "spell",
						X:      384,
						Y:      0,
						W:      48,
						H:      48,
					},
					Resource: "No Cost",
				},
			},
			wantErr: false,
		},
		{
			name:    "Test 2 - Receive invald Summoner Spells JSON",
			wantErr: true,
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					summonerSpellsJSON: []byte(`{{{{"type":"summoner","version":"9.7.1","data":{"SummonerFlash":{"id":"SummonerFlash","name":"Flash","description":"Teleports your champion a short distance toward your cursor's location.","tooltip":"Teleports your champion a short distance toward your cursor's location.","maxrank":1,"cooldown":[300],"cooldownBurn":"300","cost":[0],"costBurn":"0","datavalues":{},"effect":[null,[400],[0],[0],[0],[0],[0],[0],[0],[0],[0]],"effectBurn":[null,"400","0","0","0","0","0","0","0","0","0"],"vars":[],"key":"4","summonerLevel":7,"modes":["CLASSIC","ODIN","TUTORIAL","ARAM","ASCENSION","FIRSTBLOOD","ASSASSINATE","URF","ARSR","DOOMBOTSTEEMO","STARGUARDIAN","PROJECT","SNOWURF","ONEFORALL","GAMEMODEX","PRACTICETOOL"],"costType":"No Cost","maxammo":"-1","range":[425],"rangeBurn":"425","image":{"full":"SummonerFlash.png","sprite":"spell0.png","group":"spell","x":288,"y":0,"w":48,"h":48},"resource":"No Cost"},"SummonerHeal":{"id":"SummonerHeal","name":"Heal","description":"Restores 90-345 Health (depending on champion level) and grants 30% Movement Speed for 1 second to you and target allied champion. This healing is halved for units recently affected by Summoner Heal.","tooltip":"Restores {{ f1 }} Health and grants 30% Movement Speed for 1 second to your champion and target allied champion. This healing is halved for units recently affected by Summoner Heal.<br /><br /><span class=\"colorFFFF00\">If this spell cannot find a target, it will cast on the most wounded allied champion in range.</span>","maxrank":1,"cooldown":[240],"cooldownBurn":"240","cost":[0],"costBurn":"0","datavalues":{},"effect":[null,[0.3],[75],[15],[0.5],[826],[0.5],[0],[0],[0],[0]],"effectBurn":[null,"0.3","75","15","0.5","826","0.5","0","0","0","0"],"vars":[{"link":"@player.level","coeff":[90,105,120,135,150,165,180,195,210,225,240,255,270,285,300,315,330,345],"key":"f1"}],"key":"7","summonerLevel":1,"modes":["CLASSIC","ODIN","TUTORIAL","ARAM","ASCENSION","FIRSTBLOOD","ASSASSINATE","URF","ARSR","DOOMBOTSTEEMO","STARGUARDIAN","PROJECT","ONEFORALL","TUTORIAL_MODULE_1","TUTORIAL_MODULE_2","GAMEMODEX","PRACTICETOOL"],"costType":"No Cost","maxammo":"-1","range":[850],"rangeBurn":"850","image":{"full":"SummonerHeal.png","sprite":"spell0.png","group":"spell","x":384,"y":0,"w":48,"h":48},"resource":"No Cost"}}}`),
				},
			},
		},
		{
			name:    "Test 3 - API Call returns error",
			wantErr: true,
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					summonerSpellsJSON: []byte(""),
					err:                fmt.Errorf("Some error"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &RiotClientV4{
				config:  tt.fields.config,
				log:     tt.fields.log,
				ddragon: tt.fields.ddragon,
			}
			gotS, err := c.SummonerSpells()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.SummonerSpells() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.SummonerSpells() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
