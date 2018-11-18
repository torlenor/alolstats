package mongobackend

import (
	"context"
	"fmt"

	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/x/bsonx"
)

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

	indexView = b.client.Database(b.config.Database).Collection("summoners").Indexes()
	_, err = indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "id", Value: bsonx.Int32(1)}},
			Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
		},
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	indexView = b.client.Database(b.config.Database).Collection("champions").Indexes()
	_, err = indexView.CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bsonx.Doc{
				{Key: "key", Value: bsonx.Int32(1)}},
			Options: bsonx.Doc{{Key: "unique", Value: bsonx.Boolean(true)}},
		},
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}
