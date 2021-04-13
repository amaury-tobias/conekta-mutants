package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Cursor struct{ cursor *mongo.Cursor }

func (c *Cursor) Close(ctx context.Context) error { return c.cursor.Close(ctx) }
func (c *Cursor) Next(ctx context.Context) bool   { return c.cursor.Next(ctx) }
func (c *Cursor) Decode(val interface{}) error    { return c.cursor.Decode(val) }
func (c *Cursor) Err() error                      { return c.cursor.Err() }

type Collection struct{ collection *mongo.Collection }

func (c *Collection) Database() MongoDatabase {
	return &Database{
		database: c.collection.Database(),
	}
}
func (c *Collection) InsertOne(
	ctx context.Context,
	document interface{},
	opts ...*options.InsertOneOptions,
) (*mongo.InsertOneResult, error) {
	return c.collection.InsertOne(ctx, document, opts...)
}
func (c *Collection) Aggregate(
	ctx context.Context,
	pipeline interface{},
	opts ...*options.AggregateOptions,
) (MongoCursor, error) {
	cursor, err := c.collection.Aggregate(ctx, pipeline, opts...)
	return &Cursor{cursor: cursor}, err
}

type Database struct {
	database *mongo.Database
}

func (d *Database) Collection(name string) MongoCollection {
	return &Collection{
		collection: d.database.Collection(name),
	}
}
func (d *Database) Client() MongoClient {
	return &Client{
		client: d.database.Client(),
	}
}

type Client struct {
	client *mongo.Client
}

func (c *Client) Connect(ctx context.Context) error {
	err := c.client.Connect(ctx)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) Ping(ctx context.Context, rp *readpref.ReadPref) error {
	err := c.client.Ping(ctx, rp)
	if err != nil {
		return err
	}
	return nil
}
func (c *Client) Database(name string) MongoDatabase {
	return &Database{
		database: c.client.Database(name),
	}
}
func (c *Client) Disconnect(ctx context.Context) error { return c.client.Disconnect(ctx) }

type MongoSession struct{}

func (s *MongoSession) NewClient(opts ...*options.ClientOptions) (MongoClient, error) {
	client, err := mongo.NewClient(opts...)
	if err != nil {
		return nil, err
	}
	return &Client{client: client}, nil
}

func NewMongoSession() Session {
	return &MongoSession{}
}
