// Package mongobackend provides a MongoDB storage backend for use with the storage package.
package mongobackend

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"
	"github.com/torlenor/alolstats/riotclient"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

// Backend represents the Memory Backend
type Backend struct {
	config config.MongoBackend
	log    *logrus.Entry
	client *mongo.Client
}

func (b *Backend) GetChampions() (riotclient.ChampionList, error) {
	return riotclient.ChampionList{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetChampionsTimeStamp() time.Time {
	return time.Now()
}

func (b *Backend) StoreChampions(championList riotclient.ChampionList) error {
	return fmt.Errorf("Not implemented")
}

func (b *Backend) GetFreeRotation() (riotclient.FreeRotation, error) {
	return riotclient.FreeRotation{}, fmt.Errorf("Not implemented")
}

func (b *Backend) GetFreeRotationTimeStamp() time.Time {
	return time.Now()
}

func (b *Backend) StoreFreeRotation(freeRotation riotclient.FreeRotation) error {
	return fmt.Errorf("Not implemented")
}

// checkCollections checks if all collections needed exist and sets the correct indices
func (b *Backend) checkCollections() error {
	b.client.Database(b.config.Database).ListCollections(context.Background(), nil)

	indexView := b.client.Database(b.config.Database).Collection("matches").Indexes()

	_, err := indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "gameid", Value: bsonx.Int32(1)},
				{Key: "gameversion", Value: bsonx.Int32(1)}},
			Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
		},
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	_, err = indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "gameid", Value: bsonx.Int32(1)}},
			Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
		},
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// NewBackend creates a new Memory Backend
func NewBackend(cfg config.MongoBackend) (*Backend, error) {
	b := &Backend{
		log:    logging.Get("MongoDB Storage Backend"),
		config: cfg,
	}
	client, err := mongo.NewClient(b.config.URL)
	if err != nil {
		b.log.Errorf("Error creating new MongoDB client: %s", err)
		return nil, err
	}
	err = client.Connect(context.Background())
	if err != nil {
		b.log.Errorf("Error connecting to MongoDB: %s", err)
		return nil, err
	}
	b.client = client

	err = b.checkCollections()
	if err != nil {
		b.log.Errorf("Error checking collections: %s", err)
		return nil, err
	}

	return b, nil
}
