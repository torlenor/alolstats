package mongobackend

import (
	"context"
	"fmt"
	"time"

	"git.abyle.org/hps/alolstats/riotclient"
	"github.com/mongodb/mongo-go-driver/bson"
)

// GetFreeRotation gets the stored free champions rotation
func (b *Backend) GetFreeRotation() (*riotclient.FreeRotation, error) {
	c := b.client.Database(b.config.Database).Collection("freerotation")

	cur, err := c.Find(
		context.Background(),
		bson.D{{}},
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	freeRotation := riotclient.FreeRotation{}

	for cur.Next(nil) {
		err := cur.Decode(&freeRotation)
		if err != nil {
			b.log.Errorln("Decode error:", err)
			return nil, fmt.Errorf("Decode error: %s", err)
		}
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &freeRotation, nil
}

// GetFreeRotationTimeStamp gets the timestamp of the stored free champions rotation
func (b *Backend) GetFreeRotationTimeStamp() time.Time {
	freeRotation, err := b.GetFreeRotation()
	if err != nil {
		b.log.Errorf("Error getting Free Rotation for TimeStamp")
		return time.Time{}
	}

	return freeRotation.Timestamp
}

// StoreFreeRotation stores a new free champions rotation list
func (b *Backend) StoreFreeRotation(data *riotclient.FreeRotation) error {
	b.log.Debugf("Storing Free Rotation in storage")

	c := b.client.Database(b.config.Database).Collection("freerotation")

	// Make sure we clean possible old entry first
	_, err := c.DeleteMany(context.Background(), bson.D{})
	if err != nil {
		b.log.Debugf("%d", err)
	}

	_, err = c.InsertOne(context.Background(), data)
	if err != nil {
		return err
	}

	return nil
}
