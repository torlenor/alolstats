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

func TestRiotClientV4_Champions(t *testing.T) {
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
		wantS   riotclient.ChampionsList
		wantErr bool
	}{
		{
			name: "Test 1 - Receive valid Champions JSON",
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					championsJSON: []byte(`{"type":"champion","format":"standAloneComplex","version":"8.24.1","data":{"Aatrox":{"version":"8.24.1","id":"Aatrox","key":"266","name":"Aatrox","title":"the Darkin Blade","blurb":"Once honored defenders of Shurima against the Void, Aatrox and his brethren would eventually become an even greater threat to Runeterra, and were defeated only by cunning mortal sorcery. But after centuries of imprisonment, Aatrox was the first to find...","info":{"attack":8,"defense":4,"magic":3,"difficulty":4},"image":{"full":"Aatrox.png","sprite":"champion0.png","group":"champion","x":0,"y":0,"w":48,"h":48},"tags":["Fighter","Tank"],"partype":"Blood Well","stats":{"hp":580,"hpperlevel":80,"mp":0,"mpperlevel":0,"movespeed":345,"armor":33,"armorperlevel":3.25,"spellblock":32.1,"spellblockperlevel":1.25,"attackrange":175,"hpregen":5,"hpregenperlevel":0.25,"mpregen":0,"mpregenperlevel":0,"crit":0,"critperlevel":0,"attackdamage":60,"attackdamageperlevel":5,"attackspeedperlevel":2.5,"attackspeed":0.651}},"Ahri":{"version":"8.24.1","id":"Ahri","key":"103","name":"Ahri","title":"the Nine-Tailed Fox","blurb":"Innately connected to the latent power of Runeterra, Ahri is a vastaya who can reshape magic into orbs of raw energy. She revels in toying with her prey by manipulating their emotions before devouring their life essence. Despite her predatory nature...","info":{"attack":3,"defense":4,"magic":8,"difficulty":5},"image":{"full":"Ahri.png","sprite":"champion0.png","group":"champion","x":48,"y":0,"w":48,"h":48},"tags":["Mage","Assassin"],"partype":"Mana","stats":{"hp":526,"hpperlevel":92,"mp":418,"mpperlevel":25,"movespeed":330,"armor":20.88,"armorperlevel":3.5,"spellblock":30,"spellblockperlevel":0.5,"attackrange":550,"hpregen":6.5,"hpregenperlevel":0.6,"mpregen":8,"mpregenperlevel":0.8,"crit":0,"critperlevel":0,"attackdamage":53.04,"attackdamageperlevel":3,"attackspeedperlevel":2,"attackspeed":0.668}}}}`),
				},
			},
			wantS: riotclient.ChampionsList{
				"Aatrox": {
					Timestamp: now(),

					Version: "8.24.1",
					ID:      "Aatrox",
					Key:     "266",
					Name:    "Aatrox",
					Title:   "the Darkin Blade",
					Blurb:   "Once honored defenders of Shurima against the Void, Aatrox and his brethren would eventually become an even greater threat to Runeterra, and were defeated only by cunning mortal sorcery. But after centuries of imprisonment, Aatrox was the first to find...",
					Info: riotclient.ChampionInfo{
						Attack:     8,
						Defense:    4,
						Magic:      3,
						Difficulty: 4,
					},
					Image: riotclient.ChampionImage{
						Full:   "Aatrox.png",
						Sprite: "champion0.png",
						Group:  "champion",
						X:      0,
						Y:      0,
						W:      48,
						H:      48,
					},
					Tags: []string{
						"Fighter",
						"Tank",
					},
					Partype: "Blood Well",
					Stats: riotclient.ChampionStats{
						Hp:                   580,
						Hpperlevel:           80,
						Mp:                   0,
						Mpperlevel:           0,
						Movespeed:            345,
						Armor:                33,
						Armorperlevel:        3.25,
						Spellblock:           32.1,
						Spellblockperlevel:   1.25,
						Attackrange:          175,
						Hpregen:              5,
						Hpregenperlevel:      0.25,
						Mpregen:              0,
						Mpregenperlevel:      0,
						Crit:                 0,
						Critperlevel:         0,
						Attackdamage:         60,
						Attackdamageperlevel: 5,
						Attackspeedperlevel:  2.5,
						Attackspeed:          0.651,
					},
				},
				"Ahri": {
					Timestamp: now(),

					Version: "8.24.1",
					ID:      "Ahri",
					Key:     "103",
					Name:    "Ahri",
					Title:   "the Nine-Tailed Fox",
					Blurb:   "Innately connected to the latent power of Runeterra, Ahri is a vastaya who can reshape magic into orbs of raw energy. She revels in toying with her prey by manipulating their emotions before devouring their life essence. Despite her predatory nature...",
					Info: riotclient.ChampionInfo{
						Attack:     3,
						Defense:    4,
						Magic:      8,
						Difficulty: 5,
					},
					Image: riotclient.ChampionImage{
						Full:   "Ahri.png",
						Sprite: "champion0.png",
						Group:  "champion",
						X:      48,
						Y:      0,
						W:      48,
						H:      48,
					},
					Tags: []string{
						"Mage",
						"Assassin",
					},
					Partype: "Mana",
					Stats: riotclient.ChampionStats{
						Hp:                   526,
						Hpperlevel:           92,
						Mp:                   418,
						Mpperlevel:           25,
						Movespeed:            330,
						Armor:                20.88,
						Armorperlevel:        3.5,
						Spellblock:           30,
						Spellblockperlevel:   0.5,
						Attackrange:          550,
						Hpregen:              6.5,
						Hpregenperlevel:      0.6,
						Mpregen:              8,
						Mpregenperlevel:      0.8,
						Crit:                 0,
						Critperlevel:         0,
						Attackdamage:         53.04,
						Attackdamageperlevel: 3,
						Attackspeedperlevel:  2,
						Attackspeed:          0.668,
					},
				},
			},
			wantErr: false,
		},
		{
			name:    "Test 2 - Receive invald Champions JSON",
			wantErr: true,
			fields: fields{
				config: config.RiotClient{
					APIVersion: "v4",
					Region:     "euw1",
				},
				log: logging.Get("RiotClientV4"),
				ddragon: &MockRiotClientDD{
					championsJSON: []byte(`{{{{"type":"champion","format":"standAloneComplex","version":"8.24.1","data":{"Aatrox":{"version":"8.24.1","id":"Aatrox","key":"266","name":"Aatrox","title":"the Darkin Blade","blurb":"Once honored defenders of Shurima against the Void, Aatrox and his brethren would eventually become an even greater threat to Runeterra, and were defeated only by cunning mortal sorcery. But after centuries of imprisonment, Aatrox was the first to find...","info":{"attack":8,"defense":4,"magic":3,"difficulty":4},"image":{"full":"Aatrox.png","sprite":"champion0.png","group":"champion","x":0,"y":0,"w":48,"h":48},"tags":["Fighter","Tank"],"partype":"Blood Well","stats":{"hp":580,"hpperlevel":80,"mp":0,"mpperlevel":0,"movespeed":345,"armor":33,"armorperlevel":3.25,"spellblock":32.1,"spellblockperlevel":1.25,"attackrange":175,"hpregen":5,"hpregenperlevel":0.25,"mpregen":0,"mpregenperlevel":0,"crit":0,"critperlevel":0,"attackdamage":60,"attackdamageperlevel":5,"attackspeedperlevel":2.5,"attackspeed":0.651}},"Ahri":{"version":"8.24.1","id":"Ahri","key":"103","name":"Ahri","title":"the Nine-Tailed Fox","blurb":"Innately connected to the latent power of Runeterra, Ahri is a vastaya who can reshape magic into orbs of raw energy. She revels in toying with her prey by manipulating their emotions before devouring their life essence. Despite her predatory nature...","info":{"attack":3,"defense":4,"magic":8,"difficulty":5},"image":{"full":"Ahri.png","sprite":"champion0.png","group":"champion","x":48,"y":0,"w":48,"h":48},"tags":["Mage","Assassin"],"partype":"Mana","stats":{"hp":526,"hpperlevel":92,"mp":418,"mpperlevel":25,"movespeed":330,"armor":20.88,"armorperlevel":3.5,"spellblock":30,"spellblockperlevel":0.5,"attackrange":550,"hpregen":6.5,"hpregenperlevel":0.6,"mpregen":8,"mpregenperlevel":0.8,"crit":0,"critperlevel":0,"attackdamage":53.04,"attackdamageperlevel":3,"attackspeedperlevel":2,"attackspeed":0.668}}}}`),
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
					championsJSON: []byte(""),
					err:           fmt.Errorf("Some error"),
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
			gotS, err := c.Champions()
			if (err != nil) != tt.wantErr {
				t.Errorf("RiotClientV4.Champions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("RiotClientV4.Champions() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
