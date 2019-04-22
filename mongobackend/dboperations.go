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

// checkSummonerSpells checks the summonerspells collection and sets the correct indices
func (b *Backend) checkSummonerSpells() error {

	err := b.createIndex("summonerspells", mongo.IndexModel{
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
			{Key: "platformid", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex("matches", mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "mapid", Value: bsonx.Int32(1)},
			{Key: "queueid", Value: bsonx.Int32(1)},
		},
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

// checkChampionStats checks the championstats collection and sets the correct indices
func (b *Backend) checkChampionStats() error {
	collection := "championstats"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championkey", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championid", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkChampionStatsSummary checks the championstatssummary collection and sets the correct indices
func (b *Backend) checkChampionStatsSummary() error {
	collection := "championstatssummary"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkItemStats checks the itemstats collection and sets the correct indices
func (b *Backend) checkItemStats() error {
	collection := "itemstats"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championkey", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championid", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkSummonerSpellsStats checks the summonerspellsstats collection and sets the correct indices
func (b *Backend) checkSummonerSpellsStats() error {
	collection := "summonerspellsstats"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championkey", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championid", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkRunesReforgedStats checks the runesreforgedstats collection and sets the correct indices
func (b *Backend) checkRunesReforgedStats() error {
	collection := "runesreforgedstats"
	err := b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championkey", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
		Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	err = b.createIndex(collection, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: "championid", Value: bsonx.Int32(1)},
			{Key: "gameversion", Value: bsonx.Int32(1)},
			{Key: "tier", Value: bsonx.Int32(1)},
			{Key: "queue", Value: bsonx.Int32(1)},
		},
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

	err = b.checkChampionStats()
	if err != nil {
		return err
	}

	err = b.checkChampionStatsSummary()
	if err != nil {
		return err
	}

	err = b.checkItemStats()
	if err != nil {
		return err
	}

	err = b.checkSummonerSpellsStats()
	if err != nil {
		return err
	}

	err = b.checkRunesReforgedStats()
	if err != nil {
		return err
	}

	err = b.checkSummonerSpells()
	if err != nil {
		return err
	}

	return nil
}
