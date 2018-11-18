// Package mongobackend provides a MongoDB storage backend for use with the storage package.
package mongobackend

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/alolstats/config"
	"github.com/torlenor/alolstats/logging"

	"github.com/mongodb/mongo-go-driver/mongo"
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
