package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/storage"
)

// GetChampionStatsByChampionIDGameVersionTierQueue returns all stats specific to a certain game version, champion id and tier and queue
func (b *Backend) GetChampionStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue string) (*storage.ChampionStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("championstats")

	query := bson.D{
		{Key: "championid", Value: championID},
		{Key: "gameversion", Value: gameVersion},
		{Key: "tier", Value: tier},
		{Key: "queue", Value: queue},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Champion Stats found for Champion ID %s, GameVersion %s and Tier %s", championID, gameVersion, tier)
	}

	stat := storage.ChampionStatsStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Champion Stats for Champion ID %s, GameVersion %s, Tier %s and Queue %s: %s", championID, gameVersion, tier, queue, err)
	}

	return &stat, nil
}

// StoreChampionStats stores new champion stats in storage
func (b *Backend) StoreChampionStats(data *storage.ChampionStatsStorage) error {
	c := b.client.Database(b.config.Database).Collection("championstats")

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

// GetChampionStatsSummaryByGameVersionTierQueue returns the stats summary for a specific game version, queue and tier
func (b *Backend) GetChampionStatsSummaryByGameVersionTierQueue(gameVersion, tier, queue string) (*storage.ChampionStatsSummaryStorage, error) {
	c := b.client.Database(b.config.Database).Collection("championstatssummary")

	query := bson.D{
		{Key: "gameversion", Value: gameVersion},
		{Key: "tier", Value: tier},
		{Key: "queue", Value: queue},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Champion Stats Summary found for GameVersion %s, Queue %s and Tier %s", gameVersion, queue, tier)
	}

	stat := storage.ChampionStatsSummaryStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Champion Stats Summary for GameVersion %s, Queue %s and Tier %s: %s", gameVersion, queue, tier, err)
	}

	return &stat, nil
}

// StoreChampionStatsSummary stores a Champion Statistics Summary in the db
func (b *Backend) StoreChampionStatsSummary(data *storage.ChampionStatsSummaryStorage) error {
	c := b.client.Database(b.config.Database).Collection("championstatssummary")

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	query := bson.D{
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
