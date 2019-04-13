package mongobackend

import (
	"context"
	"fmt"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/torlenor/alolstats/riotclient"
)

// GetSummonerSpells gets the summoner spells list from storage
func (b *Backend) GetSummonerSpells() (*riotclient.SummonerSpellsList, error) {

	c := b.client.Database(b.config.Database).Collection("summonerspells")

	cur, err := c.Find(
		context.Background(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	summonerSpellsList := make(riotclient.SummonerSpellsList)

	for cur.Next(nil) {
		summonerSpell := riotclient.SummonerSpell{}
		err := cur.Decode(&summonerSpell)
		if err != nil {
			b.log.Warnln("Decode error ", err)
			continue
		}
		summonerSpellsList[summonerSpell.ID] = summonerSpell
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &summonerSpellsList, nil
}

// GetSummonerSpellsTimeStamp gets the timestamp of the stored champions list
func (b *Backend) GetSummonerSpellsTimeStamp() time.Time {
	summonerSpellsList, err := b.GetSummonerSpells()
	if err != nil {
		b.log.Errorf("Error getting Summoner Spells for TimeStamp")
		return time.Time{}
	}

	if len(*summonerSpellsList) == 0 {
		return time.Time{}
	}

	// Find oldest summoner spell time
	oldest := time.Now()
	for _, summonerSpell := range *summonerSpellsList {
		if oldest.Sub(summonerSpell.Timestamp) > 0 {
			oldest = summonerSpell.Timestamp
		}
	}

	return oldest
}

// StoreSummonerSpells stores a new champions list
func (b *Backend) StoreSummonerSpells(summnerSpellsList *riotclient.SummonerSpellsList) error {
	b.log.Debugf("Storing Summoner Spells in storage")

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	hadErrors := false
	for _, summonerSpell := range *summnerSpellsList {
		query := bson.D{{Key: "key", Value: summonerSpell.Key}}
		update := bson.D{{Key: "$set", Value: summonerSpell}}

		c := b.client.Database(b.config.Database).Collection("summonerspells")
		_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
		if err != nil {
			b.log.Warnf("Error saving Summoner Spells in DB: %s", err)
			hadErrors = true
		}
	}

	if hadErrors {
		return fmt.Errorf("There were errors saving the Summoner Spells in DB")
	}

	return nil
}
