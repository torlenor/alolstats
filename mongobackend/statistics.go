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
			Key: "gameversion", Value: gameVersion,
		},
		{
			Key: "championid", Value: championID,
		},
		{
			Key: "tier", Value: tier,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("No Champion Stats found for GameVersion %s, Champion ID %s and Tier %s: %s", gameVersion, championID, tier, err)
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

	return nil, fmt.Errorf("Could not find statistics for Champion %s, Game Version %s and tier %s", championID, gameVersion, tier)
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

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("No Champion Stats Summary found for GameVersion %s and Tier %s: %s", gameVersion, tier, err)
	}

	defer cur.Close(context.Background())

	stats := []storage.ChampionStatsSummaryStorage{}

	for cur.Next(nil) {
		stat := storage.ChampionStatsSummaryStorage{}
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
		return nil, fmt.Errorf("Found one than more Champion Stats Summary (namely %d) for Game Version %s and Tier %s in storage backend", len(stats), gameVersion, tier)
	}

	return nil, fmt.Errorf("Could not find statistics summary for Game Version %s and tier %s", gameVersion, tier)
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
