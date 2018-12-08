package mongobackend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"

	"github.com/torlenor/alolstats/riotclient"
)

func (b *Backend) summonerQuery(query *bson.D) (riotclient.Summoner, error) {
	c := b.client.Database(b.config.Database).Collection("summoners")

	cur, err := c.Find(
		context.Background(),
		query,
	)
	if err != nil {
		return riotclient.Summoner{}, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var summoners []riotclient.Summoner

	for cur.Next(nil) {
		summoner := riotclient.Summoner{}
		err := cur.Decode(&summoner)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		summoners = append(summoners, summoner)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	if len(summoners) == 1 {
		return summoners[0], nil
	} else if len(summoners) > 1 {
		return riotclient.Summoner{}, fmt.Errorf("Found one than more Summoner (namely %d) in storage backend", len(summoners))
	}

	return riotclient.Summoner{}, fmt.Errorf("Summoner not found in storage backend")
}

// GetSummonersCount returns the number of stored Summoners in the Backend
func (b *Backend) GetSummonersCount() (uint64, error) {
	c := b.client.Database(b.config.Database).Collection("summoners")

	summonersCount, err := c.Count(
		context.Background(),
		nil,
	)
	if err != nil {
		b.log.Errorf("Count error: %s", err)
		return 0, fmt.Errorf("Count error: %s", err)
	}

	return uint64(summonersCount), nil
}

// GetSummonerByNameTimeStamp gets the Timestamp when the data was stored for the Summoner specified by name
func (b *Backend) GetSummonerByNameTimeStamp(name string) time.Time {
	summoner, err := b.GetSummonerByName(name)
	if err != nil {
		return time.Time{}
	}
	return summoner.Timestamp
}

// GetSummonerByName retreives a summoner by name
func (b *Backend) GetSummonerByName(name string) (riotclient.Summoner, error) {
	query := &bson.D{{Key: "name", Value: strings.ToLower(name)}}
	return b.summonerQuery(query)
}

// GetSummonerBySummonerID retreives a summoner identified by its Summoner ID
func (b *Backend) GetSummonerBySummonerID(summonerID uint64) (riotclient.Summoner, error) {
	query := &bson.D{{Key: "id", Value: summonerID}}
	return b.summonerQuery(query)
}

// GetSummonerBySummonerIDTimeStamp retreives a summoners time stamp by its Summoner ID
func (b *Backend) GetSummonerBySummonerIDTimeStamp(summonerID uint64) time.Time {
	summoner, err := b.GetSummonerBySummonerID(summonerID)
	if err != nil {
		return time.Time{}
	}
	return summoner.Timestamp
}

// GetSummonerByAccountID retreives a summoner identified by its Account ID
func (b *Backend) GetSummonerByAccountID(accountID uint64) (riotclient.Summoner, error) {
	query := &bson.D{{Key: "accountid", Value: accountID}}
	return b.summonerQuery(query)
}

// GetSummonerByAccountIDTimeStamp retreives a summoners time stamp by its Account ID
func (b *Backend) GetSummonerByAccountIDTimeStamp(accountID uint64) time.Time {
	summoner, err := b.GetSummonerByAccountID(accountID)
	if err != nil {
		return time.Time{}
	}
	return summoner.Timestamp
}

// StoreSummoner stores new Summoner data
func (b *Backend) StoreSummoner(data *riotclient.Summoner) error {
	b.log.Debugf("Storing Summoner %s in storage", data.Name)

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	query := bson.D{{Key: "accountid", Value: data.AccountID}}
	update := bson.D{{Key: "$set", Value: data}}

	c := b.client.Database(b.config.Database).Collection("summoners")
	_, err := c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}
