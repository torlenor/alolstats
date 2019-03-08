package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/storage"
)

// GetChampionStatsByChampionIDGameVersion returns all stats specific to a certain game version and champion id
func (b *Backend) GetChampionStatsByChampionIDGameVersion(championID string, gameVersion string) (*storage.ChampionStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("championstats")

	query := bson.D{
		{
			Key: "gameversion", Value: gameVersion,
		},
		{
			Key: "championid", Value: championID,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("No Champion Stats found for GameVersion %s and Champion ID %s: %s", gameVersion, championID, err)
	}

	defer cur.Close(context.Background())

	stats := []storage.ChampionStatsStorage{}

	for cur.Next(nil) {
		stat := storage.ChampionStatsStorage{}
		err := cur.Decode(&stat)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		stats = append(stats, stat)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	if len(stats) == 1 {
		return &stats[0], nil
	} else if len(stats) > 1 {
		return nil, fmt.Errorf("Found one than more Champion Stats (namely %d) in storage backend", len(stats))
	}

	return nil, fmt.Errorf("Could not find statistics for Champion %s and Game Version %s", championID, gameVersion)
}

// GetChampionStatsByChampionIDGameVersionTier returns all stats specific to a certain game version, champion id and tier
func (b *Backend) GetChampionStatsByChampionIDGameVersionTier(championID string, gameVersion string, tier string) (*storage.ChampionStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("championstats")

	query := bson.D{
		{
			Key: "championid", Value: championID,
		},
		{
			Key: "gameversion", Value: gameVersion,
		},
		{
			Key: "tier", Value: tier,
		},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Champion Stats found for Champion ID %s, GameVersion %s and Tier %s", championID, gameVersion, tier)
	}

	stat := storage.ChampionStatsStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Champion Stats for Champion ID %s, GameVersion %s and Tier %s: %s", championID, gameVersion, tier, err)
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
	}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}

// GetChampionStatsSummaryByGameVersionTier returns the stats summary for a specific game version and tier
func (b *Backend) GetChampionStatsSummaryByGameVersionTier(gameVersion string, tier string) (*storage.ChampionStatsSummaryStorage, error) {
	c := b.client.Database(b.config.Database).Collection("championstatssummary")

	query := bson.D{
		{
			Key: "gameversion", Value: gameVersion,
		},
		{
			Key: "tier", Value: tier,
		},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Champion Stats Summary found for GameVersion %s and Tier %s", gameVersion, tier)
	}

	stat := storage.ChampionStatsSummaryStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Champion Stats Summary for GameVersion %s and Tier %s: %s", gameVersion, tier, err)
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
	}
	update := bson.D{{Key: "$set", Value: data}}

	_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}
