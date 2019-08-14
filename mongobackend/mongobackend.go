// Package mongobackend provides a MongoDB storage backend for use with the storage package.
package mongobackend

import (
	"context"

	"git.abyle.org/hps/alolstats/config"
	"git.abyle.org/hps/alolstats/logging"
	"git.abyle.org/hps/alolstats/storage"
	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Backend represents the Mongo Backend
type Backend struct {
	config config.MongoBackend
	log    *logrus.Entry
	client *mongo.Client
}

// NewBackend creates a new Mongo Backend
func NewBackend(cfg config.MongoBackend) (*Backend, error) {
	b := &Backend{
		log:    logging.Get("MongoDB Storage Backend"),
		config: cfg,
	}

	clientOptions := options.Client().ApplyURI(cfg.URL)
	client, err := mongo.NewClient(clientOptions)
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

// GetStorageSummary returns stats about the stored elements in the Backend
func (b *Backend) GetStorageSummary() (storage.Summary, error) {
	summary := storage.Summary{}
	summary.NumberOfMatches, _ = b.GetMatchesCount()
	summary.NumberOfChampions, _ = b.GetChampionsCount()
	summary.NumberOfSummoners, _ = b.GetSummonersCount()

	return summary, nil
}
