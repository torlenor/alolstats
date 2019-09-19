package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.abyle.org/hps/alolstats/riotclient"
)

type storedSummonerSpell struct {
	GameVersion string
	Language    string

	ID            string
	SummonerSpell riotclient.SummonerSpell
}

// GetSummonerSpells gets the summoner spells list from storage
func (b *Backend) GetSummonerSpells(gameVersion, language string) (riotclient.SummonerSpellsList, error) {

	c := b.client.Database(b.config.Database).Collection("summonerspells")

	query := bson.D{
		{Key: "gameVersion", Value: gameVersion},
		{Key: "language", Value: language},
	}

	cur, err := c.Find(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("GetSummonerSpells find error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	defer cur.Close(context.Background())

	summonerSpellsList := make(riotclient.SummonerSpellsList)

	for cur.Next(nil) {
		summonerSpell := storedSummonerSpell{}
		err := cur.Decode(&summonerSpell)
		if err != nil {
			b.log.Warnln("GetSummonerSpells decode error for gameversion %s and language %s: %s", gameVersion, language, err)
			continue
		}
		summonerSpellsList[summonerSpell.ID] = summonerSpell.SummonerSpell
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("GetSummonerSpells cursor error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	return summonerSpellsList, nil
}

// StoreSummonerSpells stores a new champions list
func (b *Backend) StoreSummonerSpells(gameVersion, language string, summonerSpellsList riotclient.SummonerSpellsList) error {
	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	for _, summonerSpell := range summonerSpellsList {
		summonerSpellForStorage := storedSummonerSpell{
			GameVersion:   gameVersion,
			Language:      language,
			ID:            summonerSpell.ID,
			SummonerSpell: summonerSpell,
		}

		query := bson.D{
			{Key: "gameversion", Value: gameVersion},
			{Key: "language", Value: language},
			{Key: "id", Value: summonerSpellForStorage.ID},
		}
		update := bson.D{{Key: "$set", Value: summonerSpellForStorage}}

		c := b.client.Database(b.config.Database).Collection("summonerspells")
		_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
		if err != nil {
			return fmt.Errorf("Error saving Summoner Spell %s for gameversion %s and language %s in DB: %s", summonerSpellForStorage.ID, gameVersion, language, err)
		}
	}

	return nil
}
