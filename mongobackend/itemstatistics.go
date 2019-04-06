package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/storage"
)

// GetItemStatsByChampionIDGameVersion returns all item stats specific to a certain game version and champion id
func (b *Backend) GetItemStatsByChampionIDGameVersion(championID, gameVersion string) (*storage.ItemStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("itemstats")

	query := bson.D{
		{
			Key: "championid", Value: championID,
		},
		{
			Key: "gameversion", Value: gameVersion,
		},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Champion Item Stats found for Champion ID %s and GameVersion %s", championID, gameVersion)
	}

	stat := storage.ItemStatsStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Champion Item Stats for Champion ID %s and GameVersion %s: %s", championID, gameVersion, err)
	}

	return &stat, nil
}

// StoreItemStats stores new champion item stats in storage
func (b *Backend) StoreItemStats(data *storage.ItemStatsStorage) error {
	c := b.client.Database(b.config.Database).Collection("itemstats")

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	query := bson.D{
		{Key: "championid", Value: data.ChampionID},
		{Key: "gameversion", Value: data.GameVersion},
	}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}
