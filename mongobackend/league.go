package mongobackend

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"github.com/torlenor/alolstats/riotclient"
	"github.com/torlenor/alolstats/storage"
)

func (b *Backend) GetLeagueByQueue(league string, queue string) (*riotclient.LeagueListDTO, error) {
	return nil, fmt.Errorf("Not implemented")
}

func (b *Backend) GetLeagueByQueueTimeStamp(league string, queue string) (time.Time, error) {
	return time.Time{}, fmt.Errorf("Not implemented")
}

func (b *Backend) StoreLeague(*riotclient.LeagueListDTO) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) leaguesForSummonerQuery(query *bson.D) (*storage.SummonerLeagues, error) {
	c := b.client.Database(b.config.Database).Collection("summonerleagues")

	cur, err := c.Find(
		context.Background(),
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var leagues []storage.SummonerLeagues

	for cur.Next(nil) {
		league := storage.SummonerLeagues{}
		err := cur.Decode(&league)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		leagues = append(leagues, league)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	if len(leagues) == 1 {
		return &leagues[0], nil
	} else if len(leagues) > 1 {
		return nil, fmt.Errorf("Found one than more Summoner League informations (namely %d) in storage backend", len(leagues))
	}

	return nil, fmt.Errorf("Summoner League information not found in storage backend")
}

// GetLeaguesForSummoner returns the stores League information for a certain Summoner identified by its name
func (b *Backend) GetLeaguesForSummoner(summonerName string) (*storage.SummonerLeagues, error) {
	query := &bson.D{{Key: "summonername", Value: strings.ToLower(summonerName)}}
	return b.leaguesForSummonerQuery(query)
}

// GetLeaguesForSummonerTimeStamp gets the timestamp of the stored Leagues for Summoner idenfitied by name
func (b *Backend) GetLeaguesForSummonerTimeStamp(summonerName string) (time.Time, error) {
	leagues, err := b.GetLeaguesForSummoner(summonerName)
	if err != nil {
		b.log.Errorf("Error getting Leagues for Summoner %s for TimeStamp: %s", summonerName, err)
		return time.Time{}, fmt.Errorf("Error getting Leagues for Summoner %s for TimeStamp: %s", summonerName, err)
	}

	if len(leagues.LeaguePositionDTOList.LeaguePosition) == 0 {
		return time.Time{}, nil
	}

	// Find oldest champ time
	oldest := time.Now()
	for _, league := range leagues.LeaguePositionDTOList.LeaguePosition {
		if oldest.Sub(league.Timestamp) > 0 {
			oldest = league.Timestamp
		}
	}

	return oldest, nil
}

// GetLeaguesForSummonerBySummonerID returns the stores League information for a certain Summoner identified by its Summoner ID
func (b *Backend) GetLeaguesForSummonerBySummonerID(summonerID string) (*storage.SummonerLeagues, error) {
	query := &bson.D{{Key: "summonerid", Value: summonerID}}
	return b.leaguesForSummonerQuery(query)
}

// GetLeaguesForSummonerBySummonerIDTimeStamp gets the timestamp of the stored Leagues for Summoner idenfitied by ID
func (b *Backend) GetLeaguesForSummonerBySummonerIDTimeStamp(summonerID string) (time.Time, error) {
	leagues, err := b.GetLeaguesForSummonerBySummonerID(summonerID)
	if err != nil {
		b.log.Errorf("Error getting Leagues for Summoner ID %s for TimeStamp: %s", summonerID, err)
		return time.Time{}, fmt.Errorf("Error getting Leagues for Summoner ID %s for TimeStamp: %s", summonerID, err)
	}

	if len(leagues.LeaguePositionDTOList.LeaguePosition) == 0 {
		return time.Time{}, nil
	}

	// Find oldest champ time
	oldest := time.Now()
	for _, league := range leagues.LeaguePositionDTOList.LeaguePosition {
		if oldest.Sub(league.Timestamp) > 0 {
			oldest = league.Timestamp
		}
	}

	return oldest, nil
}

// StoreLeaguesForSummoner stores new Summoner League information data
func (b *Backend) StoreLeaguesForSummoner(leagues *storage.SummonerLeagues) error {
	b.log.Debugf("Storing Summoner %s in storage", leagues.SummonerName)

	c := b.client.Database(b.config.Database).Collection("summonerleagues")

	// Make sure we clean possible old entries for the same Summoner
	filter := bson.D{{Key: "summonerid", Value: leagues.SummonerID}}
	_, err := c.DeleteOne(context.Background(), filter)
	if err != nil {
		b.log.Debugf("%d", err)
	}

	upsert := true
	updateOptions := options.UpdateOptions{Upsert: &upsert}

	query := bson.D{{Key: "summonername", Value: leagues.SummonerName}}
	update := bson.D{{Key: "$set", Value: leagues}}

	_, err = c.UpdateOne(context.Background(), query, update, &updateOptions)
	if err != nil {
		return err
	}

	return nil
}
