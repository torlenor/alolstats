package mongobackend

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.abyle.org/hps/alolstats/riotclient"
)

// GetChampions gets the champions list from storage
func (b *Backend) GetChampions() (riotclient.ChampionsList, error) {

	c := b.client.Database(b.config.Database).Collection("champions")

	cur, err := c.Find(
		context.Background(),
		bson.D{{}},
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	championsList := make(riotclient.ChampionsList)

	for cur.Next(nil) {
		champion := riotclient.Champion{}
		err := cur.Decode(&champion)
		if err != nil {
			b.log.Warnln("Decode error ", err)
			continue
		}
		championsList[champion.ID] = champion
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return championsList, nil
}

// GetChampionsCount gets the number of champions from storage
func (b *Backend) GetChampionsCount() (uint64, error) {

	c := b.client.Database(b.config.Database).Collection("champions")

	championsCount, err := c.CountDocuments(
		context.Background(),
		bson.D{{}},
	)
	if err != nil {
		return 0, fmt.Errorf("Find error: %s", err)
	}

	return uint64(championsCount), nil
}

// GetChampionsTimeStamp gets the timestamp of the stored champions list
func (b *Backend) GetChampionsTimeStamp() time.Time {
	championList, err := b.GetChampions()
	if err != nil {
		b.log.Errorf("Error getting Champions for TimeStamp: %s", err)
		return time.Time{}
	}

	if len(championList) == 0 {
		return time.Time{}
	}

	// Find oldest champ time
	oldest := time.Now()
	for _, champion := range championList {
		if oldest.Sub(champion.Timestamp) > 0 {
			oldest = champion.Timestamp
		}
	}

	return oldest
}

// StoreChampions stores a new champions list
func (b *Backend) StoreChampions(championList riotclient.ChampionsList) error {
	b.log.Debugf("Storing Champions in storage")

	hadErrors := false
	for _, champion := range championList {
		query := bson.D{{Key: "key", Value: champion.Key}}
		update := bson.D{{Key: "$set", Value: champion}}

		c := b.client.Database(b.config.Database).Collection("champions")
		_, err := c.UpdateOne(context.Background(), query, update, options.Update().SetUpsert(true))
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
