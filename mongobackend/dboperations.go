package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

func (b *Backend) createIndex(collection string, indexModel mongo.IndexModel) error {
	indexView := b.client.Database(b.config.Database).Collection(collection).Indexes()
	_, err := indexView.CreateOne(
		context.Background(),
		indexModel,
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkChampions checks the champions collection and sets the correct indices
func (b *Backend) checkChampions() error {

	err := b.createIndex("champions", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "key", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkMatches checks the matches collection and sets the correct indices
func (b *Backend) checkMatches() error {
	err := b.createIndex("matches", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameid", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("matches", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("matches", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "mapid", Value: bsonx.Int32(1)},
			{Key: "queueid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(false)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("matches", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "mapid", Value: bsonx.Int32(1)},
			{Key: "queueid", Value: bsonx.Int32(1)},
			{Key: "participants.championid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(false)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkSummoners checks the summoners collection and sets the correct indices
func (b *Backend) checkSummoners() error {
	err := b.createIndex("summoners", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "summonername", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("summoners", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "accountid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("summoners", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "summonerid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("summoners", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "puuid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkSummonerLeagues checks the summonerleagues collection and sets the correct indices
func (b *Backend) checkSummonerLeagues() error {
	collection := "summonerleagues"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "summonername", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "summonerid", Value: bsonx.Int32(1)}},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkCollections checks if all collections needed exist and sets the correct indices
func (b *Backend) checkCollections() error {
	err := b.checkChampions()
	if err != nil {
		return err
	}

	err = b.checkMatches()
	if err != nil {
		return err
	}

	err = b.checkSummoners()
	if err != nil {
		return err
	}

	err = b.checkSummonerLeagues()
	if err != nil {
		return err
	}

	return nil
}
