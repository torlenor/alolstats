package mongobackend

import (
	"context"
	"fmt"

	"git.abyle.org/hps/alolstats/storage"
	"go.mongodb.org/mongo-driver/bson"
)

// GetMatchesCursorByGameVersion returns cursor to matches specific to a certain game version
func (b *Backend) GetMatchesCursorByGameVersion(gameVersion string) (storage.QueryCursor, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{{Key: "gameversion",
		Value: bson.D{
			{Key: "$regex", Value: "^" + gameVersion},
		},
	}}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error finding matches for GameVersion %s: %s", gameVersion, err)
	}

	matchCursor := MatchCursor{
		cur: cur,
		ctx: context.Background(),
	}

	return &matchCursor, nil
}

// GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs returns cursor to matches specific to a certain game version, champion id, map id and queue ids between and equal to ltequeue <= queueid <= gtequeue
func (b *Backend) GetMatchesCursorByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (storage.QueryCursor, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion},
			},
		},
		{
			Key:   "mapid",
			Value: mapID,
		},
		{
			Key: "queueid",
			Value: bson.D{
				{Key: "$lte", Value: ltequeue},
				{Key: "$gte", Value: gtequeue},
			},
		},
		{
			Key:   "participants.championid",
			Value: championID,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error finding matches for GameVersion %s, Champion ID %d, Map ID %d, Queue ID  %d <= id <= %d: %s", gameVersion, championID, mapID, gtequeue, ltequeue, err)
	}

	matchCursor := MatchCursor{
		cur: cur,
		ctx: context.Background(),
	}

	return &matchCursor, nil
}

// GetMatchesCursorByGameVersionMapBetweenQueueIDs returns cursor to matches specific to a certain game version, map id and queue ids between and equal to ltequeue <= queueid <= gtequeue
func (b *Backend) GetMatchesCursorByGameVersionMapBetweenQueueIDs(gameVersion string, mapID uint64, ltequeue uint64, gtequeue uint64) (storage.QueryCursor, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion},
			},
		},
		{
			Key:   "mapid",
			Value: mapID,
		},
		{
			Key: "queueid",
			Value: bson.D{
				{Key: "$lte", Value: ltequeue},
				{Key: "$gte", Value: gtequeue},
			},
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error finding matches for GameVersion %s, Map ID %d, Queue ID  %d <= id <= %d: %s", gameVersion, mapID, gtequeue, ltequeue, err)
	}

	matchCursor := MatchCursor{
		cur: cur,
		ctx: context.Background(),
	}

	return &matchCursor, nil
}

// GetMatchesCursorByGameVersionMapQueueID returns cursor to matches specific to a certain game version, map id and queue id
func (b *Backend) GetMatchesCursorByGameVersionMapQueueID(gameVersion string, mapID uint64, queueID uint64) (storage.QueryCursor, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion},
			},
		},
		{
			Key:   "mapid",
			Value: mapID,
		},
		{
			Key:   "queueid",
			Value: queueID,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error finding matches for GameVersion %s, Map ID %d, Queue ID  %d: %s", gameVersion, mapID, queueID, err)
	}

	matchCursor := MatchCursor{
		cur: cur,
		ctx: context.Background(),
	}

	return &matchCursor, nil
}

// GetMatchesCursorByGameVersionMajorMinorMapQueueID returns cursor to matches specific to a certain game version, map id and queue id
func (b *Backend) GetMatchesCursorByGameVersionMajorMinorMapQueueID(major int, minor int, mapID uint64, queueID uint64) (storage.QueryCursor, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key:   "gameversion_major",
			Value: major,
		},
		{
			Key:   "gameversion_minor",
			Value: minor,
		},
		{
			Key:   "mapid",
			Value: mapID,
		},
		{
			Key:   "queueid",
			Value: queueID,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("Error finding matches for GameVersion %d%d, Map ID %d, Queue ID  %d: %s", major, minor, mapID, queueID, err)
	}

	matchCursor := MatchCursor{
		cur: cur,
		ctx: context.Background(),
	}

	return &matchCursor, nil
}
