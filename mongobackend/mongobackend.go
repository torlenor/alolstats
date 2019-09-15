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
	b.client = client

	return b, nil
}

// Connect initializes the MongoDBBackend Client by using the appropriate Mongo functions.
// In addition it checks its collections and creates indices if necessary.
// This method must be called before the Mongo Backend can be used.
func (b *Backend) Connect() error {
	err := b.client.Connect(context.Background())
	if err != nil {
		b.log.Errorf("Error connecting to MongoDB: %s", err)
		return err
	}

	err = b.checkCollections()
	if err != nil {
		b.log.Errorf("Error checking collections: %s", err)
		return err
	}

	return nil
}

// GetStorageSummary returns stats about the stored elements in the Backend
func (b *Backend) GetStorageSummary() (storage.Summary, error) {
	summary := storage.Summary{}
	summary.NumberOfMatches, _ = b.GetMatchesCount()
	summary.NumberOfChampions, _ = b.GetChampionsCount()
	summary.NumberOfSummoners, _ = b.GetSummonersCount()

	return summary, nil
}
