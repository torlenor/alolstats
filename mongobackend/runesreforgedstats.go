package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/storage"
)

// GetRunesReforgedStatsByChampionIDGameVersionTierQueue returns all stats specific to a certain game version, champion id and tier and queue
func (b *Backend) GetRunesReforgedStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue string) (*storage.RunesReforgedStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("runesreforgedstats")

	query := bson.D{
		{Key: "championid", Value: championID},
		{Key: "gameversion", Value: gameVersion},
		{Key: "tier", Value: tier},
		{Key: "queue", Value: queue},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Runes Reforged Stats found for Champion ID %s, GameVersion %s and Tier %s", championID, gameVersion, tier)
	}

	stat := storage.RunesReforgedStatsStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Runes Reforged Stats for Champion ID %s, GameVersion %s, Tier %s and Queue %s: %s", championID, gameVersion, tier, queue, err)
	}

	return &stat, nil
}

// StoreRunesReforgedStats stores new summonerspells stats in storage
func (b *Backend) StoreRunesReforgedStats(data *storage.RunesReforgedStatsStorage) error {
	c := b.client.Database(b.config.Database).Collection("runesreforgedstats")

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	query := bson.D{
		{Key: "championid", Value: data.ChampionID},
		{Key: "gameversion", Value: data.GameVersion},
		{Key: "tier", Value: data.Tier},
		{Key: "queue", Value: data.Queue},
	}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}
