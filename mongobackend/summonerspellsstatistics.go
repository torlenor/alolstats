package mongobackend

import (
	"context"
	"fmt"

	"git.abyle.org/hps/alolstats/storage"
	"github.com/mongodb/mongo-go-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetSummonerSpellsStatsByChampionIDGameVersionTierQueue returns all stats specific to a certain game version, champion id and tier and queue
func (b *Backend) GetSummonerSpellsStatsByChampionIDGameVersionTierQueue(championID, gameVersion, tier, queue string) (*storage.SummonerSpellsStatsStorage, error) {
	c := b.client.Database(b.config.Database).Collection("summonerspellsstats")

	query := bson.D{
		{Key: "championid", Value: championID},
		{Key: "gameversion", Value: gameVersion},
		{Key: "tier", Value: tier},
		{Key: "queue", Value: queue},
	}

	doc := c.FindOne(
		context.Background(), query)
	if doc == nil {
		return nil, fmt.Errorf("No Summoner Spells Stats found for Champion ID %s, GameVersion %s, Tier %s and Queue %s", championID, gameVersion, tier, queue)
	}

	stat := storage.SummonerSpellsStatsStorage{}
	err := doc.Decode(&stat)
	if err != nil {
		return nil, fmt.Errorf("Decode error when trying to Decode Summoner Spells Stats for Champion ID %s, GameVersion %s, Tier %s and Queue %s: %s", championID, gameVersion, tier, queue, err)
	}

	return &stat, nil
}

// StoreSummonerSpellsStats stores new summonerspells stats in storage
func (b *Backend) StoreSummonerSpellsStats(data *storage.SummonerSpellsStatsStorage) error {
	c := b.client.Database(b.config.Database).Collection("summonerspellsstats")

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
