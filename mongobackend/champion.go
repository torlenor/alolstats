package mongobackend

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions gets the champions list from storage
func (b *Backend) GetChampions() (riotclient.ChampionList, error) {

	c := b.client.Database(b.config.Database).Collection("champions")

	cur, err := c.Find(
		context.Background(),
		nil,
	)
	if err != nil {
		return riotclient.ChampionList{}, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var championsList riotclient.ChampionList
	championsList.Champions = make(map[string]riotclient.Champion)

	for cur.Next(nil) {
		champion := riotclient.Champion{}
		err := cur.Decode(&champion)
		if err != nil {
			b.log.Warnln("Decode error ", err)
			continue
		}
		championsList.Champions[champion.ID] = champion
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return championsList, nil
}

// GetChampionsTimeStamp gets the timestamp of the stored champions list
func (b *Backend) GetChampionsTimeStamp() time.Time {
	championList, err := b.GetChampions()
	if err != nil {
		b.log.Errorf("Error getting Champions for TimeStamp")
		return time.Time{}
	}

	// Find oldest champ time
	oldest := time.Now()
	for _, champion := range championList.Champions {
		if oldest.Sub(champion.Timestamp) > 0 {
			oldest = champion.Timestamp
		}
	}

	return oldest
}

// StoreChampions stores a new champions list
func (b *Backend) StoreChampions(championList riotclient.ChampionList) error {
	b.log.Debugf("Storing Champions in storage")

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	hadErrors := false
	for _, champion := range championList.Champions {
		query := bson.D{{Key: "key", Value: champion.Key}}
		update := bson.D{{Key: "$set", Value: champion}}

		c := b.client.Database(b.config.Database).Collection("champions")
		_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
		if err != nil {
			b.log.Warnf("Error saving champion in DB: %s", err)
			hadErrors = true
		}
	}

	if hadErrors {
		return fmt.Errorf("There were errors saving the cmapions in DB")
	}

	return nil
}
