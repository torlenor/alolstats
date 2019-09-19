package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.abyle.org/hps/alolstats/riotclient"
)

type storedRunesReforged struct {
	GameVersion string
	Language    string

	ID            int
	RunesReforged riotclient.RunesReforgedSet
}

// GetRunesReforged gets the summoner spells list from storage
func (b *Backend) GetRunesReforged(gameVersion, language string) (riotclient.RunesReforgedList, error) {

	c := b.client.Database(b.config.Database).Collection("runesreforged")

	query := bson.D{
		{Key: "gameVersion", Value: gameVersion},
		{Key: "language", Value: language},
	}

	cur, err := c.Find(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("GetRunesReforged find error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	defer cur.Close(context.Background())

	rfList := make(riotclient.RunesReforgedList)

	for cur.Next(nil) {
		runerf := storedRunesReforged{}
		err := cur.Decode(&runerf)
		if err != nil {
			b.log.Warnln("GetRunesReforged decode error for gameversion %s and language %s: %s", gameVersion, language, err)
			continue
		}
		rfList[runerf.ID] = runerf.RunesReforged
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("GetRunesReforged cursor error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	return rfList, nil
}

// StoreRunesReforged stores a new champions list
func (b *Backend) StoreRunesReforged(gameVersion, language string, rfList riotclient.RunesReforgedList) error {
	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	for _, rfset := range rfList {
		runesForStorage := storedRunesReforged{
			GameVersion:   gameVersion,
			Language:      language,
			ID:            rfset.ID,
			RunesReforged: rfset,
		}

		query := bson.D{
			{Key: "gameVersion", Value: gameVersion},
			{Key: "language", Value: language},
			{Key: "id", Value: runesForStorage.ID},
		}
		update := bson.D{{Key: "$set", Value: runesForStorage}}

		c := b.client.Database(b.config.Database).Collection("runesreforged")
		_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
		if err != nil {
			return fmt.Errorf("Error saving RunesReforged %d for gameversion %s and language %s in DB: %s", runesForStorage.ID, gameVersion, language, err)
		}
	}

	return nil
}
