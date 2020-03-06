package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.abyle.org/hps/alolstats/riotclient"
)

type storedItem struct {
	GameVersion string
	Language    string

	Key  uint16
	Item riotclient.Item
}

// GetItems gets the summoner spells list from storage
func (b *Backend) GetItems(gameVersion, language string) (riotclient.ItemList, error) {

	c := b.client.Database(b.config.Database).Collection("items")

	query := bson.D{
		{Key: "gameVersion", Value: gameVersion},
		{Key: "language", Value: language},
	}

	cur, err := c.Find(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("GetItems find error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	defer cur.Close(context.Background())

	itemList := make(riotclient.ItemList)

	for cur.Next(nil) {
		item := storedItem{}
		err := cur.Decode(&item)
		if err != nil {
			b.log.Warnln("GetItems decode error for gameversion %s and language %s: %s", gameVersion, language, err)
			continue
		}
		itemList[item.Key] = item.Item
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("GetItems cursor error for gameversion %s and language %s: %s", gameVersion, language, err)
	}

	return itemList, nil
}

// StoreItems stores a new champions list
func (b *Backend) StoreItems(gameVersion, language string, itemList riotclient.ItemList) error {
	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	for _, item := range itemList {
		itemForStorage := storedItem{
			GameVersion: gameVersion,
			Language:    language,
			Key:         item.Key,
			Item:        item,
		}

		query := bson.D{
			{Key: "gameVersion", Value: gameVersion},
			{Key: "language", Value: language},
			{Key: "id", Value: itemForStorage.Key},
		}
		update := bson.D{{Key: "$set", Value: itemForStorage}}

		c := b.client.Database(b.config.Database).Collection("items")
		_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
		if err != nil {
			return fmt.Errorf("Error saving Item %d for gameversion %s and language %s in DB: %s", itemForStorage.Key, gameVersion, language, err)
		}
	}

	return nil
}
