package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/torlenor/alolstats/riotclient"
)

// GetMatch retreives match data for given id
func (b *Backend) GetMatch(id uint64) (*riotclient.MatchDTO, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	cur, err := c.Find(
		context.Background(),
		bson.D{{Key: "gameid", Value: id}},
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var matches []riotclient.MatchDTO

	for cur.Next(nil) {
		match := riotclient.MatchDTO{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches = append(matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	if len(matches) == 1 {
		return &matches[0], nil
	} else if len(matches) > 1 {
		return nil, fmt.Errorf("Found one than more Match (namely %d) with id=%d in storage backend", len(matches), id)
	}

	return nil, fmt.Errorf("Match with id=%d not found in storage backend", id)
}

// GetMatchesCount returns the number of stored Matches in the Backend
func (b *Backend) GetMatchesCount() (uint64, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	matchesCount, err := c.Count(
		context.Background(),
		nil,
	)
	if err != nil {
		return uint64(0), fmt.Errorf("Find error: %s", err)
	}

	return uint64(matchesCount), nil
}

// StoreMatch stores new match data
func (b *Backend) StoreMatch(data *riotclient.MatchDTO) error {
	b.log.Debugf("Storing Match id=%d in storage", data.GameID)

	c := b.client.Database(b.config.Database).Collection("matches")
	_, err := c.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

// GetMatchesByGameVersion returns all matches specific to a certain game version from Storage
func (b *Backend) GetMatchesByGameVersion(gameVersion string) (*riotclient.Matches, error) {

	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{{Key: "gameversion",
		Value: bson.D{
			{Key: "$regex", Value: "^" + gameVersion + ""},
		},
	}}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("No match found for GameVersion %s: %s", gameVersion, err)
	}

	defer cur.Close(context.Background())

	matches := riotclient.Matches{}

	for cur.Next(nil) {
		match := riotclient.MatchDTO{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches.Matches = append(matches.Matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &matches, nil
}

// GetMatchesByGameVersionAndChampionID returns all matches specific to a certain game version and champion id
func (b *Backend) GetMatchesByGameVersionAndChampionID(gameVersion string, championID uint64) (*riotclient.Matches, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion + ""},
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
		return nil, fmt.Errorf("No match found for GameVersion %s and Champion ID %d: %s", gameVersion, championID, err)
	}

	defer cur.Close(context.Background())

	matches := riotclient.Matches{}

	for cur.Next(nil) {
		match := riotclient.MatchDTO{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches.Matches = append(matches.Matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &matches, nil
}

// GetMatchesByGameVersionChampionIDMapQueue returns all matches specific to a certain game version, champion id, map id and queue id
func (b *Backend) GetMatchesByGameVersionChampionIDMapQueue(gameVersion string, championID uint64, mapID uint64, queueID uint64) (*riotclient.Matches, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion + ""},
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
		{
			Key:   "participants.championid",
			Value: championID,
		},
	}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return nil, fmt.Errorf("No match found for GameVersion %s and Champion ID %d: %s", gameVersion, championID, err)
	}

	defer cur.Close(context.Background())

	matches := riotclient.Matches{}

	for cur.Next(nil) {
		match := riotclient.MatchDTO{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches.Matches = append(matches.Matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &matches, nil
}

// GetMatchesByGameVersionChampionIDMapBetweenQueueIDs returns all matches specific to a certain game version, champion id, map id and queue ids between and equal to ltequeue <= queueid <= gtequeue
func (b *Backend) GetMatchesByGameVersionChampionIDMapBetweenQueueIDs(gameVersion string, championID uint64, mapID uint64, ltequeue uint64, gtequeue uint64) (*riotclient.Matches, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{
		{
			Key: "gameversion",
			Value: bson.D{
				{Key: "$regex", Value: "^" + gameVersion + ""},
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
		return nil, fmt.Errorf("No match found for GameVersion %s and Champion ID %d: %s", gameVersion, championID, err)
	}

	defer cur.Close(context.Background())

	matches := riotclient.Matches{}

	for cur.Next(nil) {
		match := riotclient.MatchDTO{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches.Matches = append(matches.Matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &matches, nil
}
