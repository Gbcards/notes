package platform

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/fx"
)

const (
	defaultMongoURI = "mongodb://localhost:27017"
	mongoDatabase   = "notes_service"
	notesCollection = "notes"
)

func mongoURI() string {
	if uri := os.Getenv("MONGO_URI"); uri != "" {
		return uri
	}
	return defaultMongoURI
}

func NewMongoClient(lc fx.Lifecycle) (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI()))
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return client.Disconnect(ctx)
		},
	})

	return client, nil
}

func NewNoteCollection(client *mongo.Client) *mongo.Collection {
	return client.Database(mongoDatabase).Collection(notesCollection)
}
