package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/torlenor/alolstats/riotclient"
)

// GetMatch retreives match data for given id
func (b *Backend) GetMatch(id uint64) (riotclient.Match, error) {
	c := b.client.Database(b.config.Database).Collection("matches")

	cur, err := c.Find(
		context.Background(),
		bson.D{{Key: "gameid", Value: id}},
	)
	if err != nil {
		return riotclient.Match{}, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var matches []riotclient.Match

	for cur.Next(nil) {
		match := riotclient.Match{}
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
		return matches[0], nil
	} else if len(matches) > 1 {
		return riotclient.Match{}, fmt.Errorf("Found one than more Match (namely %d) with id=%d in storage backend", len(matches), id)
	}

	return riotclient.Match{}, fmt.Errorf("Match with id=%d not found in storage backend", id)
}

// StoreMatch stores new match data
func (b *Backend) StoreMatch(data *riotclient.Match) error {
	b.log.Debugf("Storing Match id=%d in storage", data.GameID)

	c := b.client.Database(b.config.Database).Collection("matches")
	_, err := c.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}

// GetMatchesByGameVersion returns all matches specific to a certain game version
func (b *Backend) GetMatchesByGameVersion(gameVersion string) (riotclient.Matches, error) {

	c := b.client.Database(b.config.Database).Collection("matches")

	query := bson.D{{Key: "gameversion",
		Value: bson.D{
			{Key: "$regex", Value: "^" + gameVersion + ""},
		},
	}}

	cur, err := c.Find(
		context.Background(), query)
	if err != nil {
		return riotclient.Matches{}, fmt.Errorf("No match found for GameVersion %s: %s", gameVersion, err)
	}

	defer cur.Close(context.Background())

	matches := riotclient.Matches{}

	for cur.Next(nil) {
		match := riotclient.Match{}
		err := cur.Decode(&match)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		matches.Matches = append(matches.Matches, match)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return matches, nil
}
