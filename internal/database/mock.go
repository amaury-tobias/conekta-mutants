package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Session interface {
	NewClient(opts ...*options.ClientOptions) (MongoClient, error)
}
type MongoClient interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string) MongoDatabase
	Disconnect(ctx context.Context) error
}
type MongoDatabase interface {
	Client() MongoClient
	Collection(name string) MongoCollection
}
type MongoCollection interface {
	Database() MongoDatabase
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error)
	Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (MongoCursor, error)
}
type MongoCursor interface {
	Close(ctx context.Context) error
	Next(ctx context.Context) bool
	Decode(val interface{}) error
	Err() error
}

type MockCursor struct{}

func (c *MockCursor) Close(ctx context.Context) error { return nil }
func (c *MockCursor) Next(ctx context.Context) bool   { return false }
func (c *MockCursor) Decode(val interface{}) error    { return nil }
func (c *MockCursor) Err() error                      { return nil }

type MockCollection struct{}

func (c *MockCollection) Database() MongoDatabase { return &MockDatabase{} }
func (c *MockCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	return nil, nil
}
func (c *MockCollection) Aggregate(ctx context.Context, pipeline interface{}, opts ...*options.AggregateOptions) (MongoCursor, error) {
	return &MockCursor{}, nil
}

type MockDatabase struct{}

func (c *MockDatabase) Collection(name string) MongoCollection { return &MockCollection{} }
func (c *MockDatabase) Client() MongoClient                    { return &MockClient{} }

type MockClient struct{}

func (c *MockClient) Connect(ctx context.Context) error                     { return nil }
func (c *MockClient) Ping(ctx context.Context, rp *readpref.ReadPref) error { return nil }
func (c *MockClient) Database(name string) MongoDatabase                    { return &MockDatabase{} }
func (c *MockClient) Disconnect(ctx context.Context) error                  { return nil }

type MockSession struct{}

func (s *MockSession) NewClient(opts ...*options.ClientOptions) (MongoClient, error) {
	return &MockClient{}, nil
}

func NewMockSession() Session {
	return &MockSession{}
}
