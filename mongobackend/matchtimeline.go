package mongobackend

import (
	"context"
	"fmt"

	"git.abyle.org/hps/alolstats/riotclient"
	"github.com/mongodb/mongo-go-driver/bson"
)

// MatchTimeLineStoreData is used to store Match TimeLine Data and return it
type MatchTimeLineStoreData struct {
	GameID      int64
	QueueID     int
	MapID       int
	SeasonID    int
	GameVersion string
	GameMode    string
	GameType    string
	PlatformID  string
	TimeLine    *riotclient.MatchTimelineDTO
}

// GetMatchTimeLine is not implemented, yet
func (b *Backend) GetMatchTimeLine(id uint64) (*riotclient.MatchTimelineDTO, error) {
	c := b.client.Database(b.config.Database).Collection("timelines")

	cur, err := c.Find(
		context.Background(),
		bson.D{{Key: "gameid", Value: id}},
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	var timelines []MatchTimeLineStoreData

	for cur.Next(nil) {
		timeline := MatchTimeLineStoreData{}
		err := cur.Decode(&timeline)
		if err != nil {
			b.log.Warnln("Decode error ", err)
		}
		timelines = append(timelines, timeline)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	if len(timelines) == 1 {
		return timelines[0].TimeLine, nil
	} else if len(timelines) > 1 {
		return nil, fmt.Errorf("Found one than more Match TimeLine (namely %d) with id=%d in storage backend", len(timelines), id)
	}

	return nil, fmt.Errorf("Match TimeLine with id=%d not found in storage backend", id)
}

// StoreMatchTimeLine is not implemented, yet
func (b *Backend) StoreMatchTimeLine(match *riotclient.MatchDTO, data *riotclient.MatchTimelineDTO) error {
	b.log.Debugf("Storing Match TimeLine id=%d in storage", match.GameID)

	storeData := MatchTimeLineStoreData{
		GameID:      match.GameID,
		QueueID:     match.QueueID,
		MapID:       match.MapID,
		SeasonID:    match.SeasonID,
		GameVersion: match.GameVersion,
		GameMode:    match.GameMode,
		GameType:    match.GameType,
		PlatformID:  match.PlatformID,
		TimeLine:    data,
	}

	c := b.client.Database(b.config.Database).Collection("timelines")
	_, err := c.InsertOne(context.Background(), storeData)
	if err != nil {
		return err
	}

	return nil
}
