package mongobackend

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/options"
	"github.com/torlenor/alolstats/riotclient"
)

// GetChampions gets the champions list from storage
func (b *Backend) GetChampions() (riotclient.ChampionList, error) {
	return riotclient.ChampionList{}, fmt.Errorf("Not implemented")
}

// GetChampionsTimeStamp gets the timestamp of the stored champions list
func (b *Backend) GetChampionsTimeStamp() time.Time {
	return time.Now()
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
	} else {
		return nil
	}
}
