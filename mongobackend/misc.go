package mongobackend

import (
	"context"
	"fmt"

	"git.abyle.org/hps/alolstats/storage"
	"github.com/mongodb/mongo-go-driver/bson"
)

// GetKnownGameVersions gets the stored known game versions
func (b *Backend) GetKnownGameVersions() (*storage.GameVersions, error) {
	c := b.client.Database(b.config.Database).Collection("gameversions")

	cur, err := c.Find(
		context.Background(),
		bson.D{{}},
	)
	if err != nil {
		return nil, fmt.Errorf("Find error: %s", err)
	}

	defer cur.Close(context.Background())

	gameVersions := storage.GameVersions{}

	for cur.Next(nil) {
		err := cur.Decode(&gameVersions)
		if err != nil {
			b.log.Errorln("Decode error:", err)
			return nil, fmt.Errorf("Decode error: %s", err)
		}
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return &gameVersions, nil
}

// StoreKnownGameVersions stores a new list of known game versions
func (b *Backend) StoreKnownGameVersions(data *storage.GameVersions) error {
	b.log.Debugf("Storing known Game Versions in storage")

	c := b.client.Database(b.config.Database).Collection("gameversions")

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
