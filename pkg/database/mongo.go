package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewDatabase(url string) (*Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(url))
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	// Extract database name from connection string (or set default)
	dbName := "ranking_video" // You can extract this from URL or pass separately
	database := &Database{
		client:   client,
		database: client.Database(dbName),
	}

	return database, nil
}

// GetCollection returns a MongoDB collection
func (d *Database) GetCollection(name string) *mongo.Collection {
	return d.database.Collection(name)
}

// GetDatabase returns the MongoDB database
func (d *Database) GetDatabase() *mongo.Database {
	return d.database
}

// Close closes the database connection
func (d *Database) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.client.Disconnect(ctx)
}
