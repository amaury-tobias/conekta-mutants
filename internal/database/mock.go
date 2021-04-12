package database

import (
	"context"

	"github.com/amaury-tobias/conekta-mutants/internal/models"
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
	GetStats() (*models.Stats, error)
	SaveHuman(*models.HumanModel) error
}

type MockCollection struct{}

func (c *MockCollection) GetStats() (*models.Stats, error) {
	return &models.Stats{
		CountMutantDNA: 5,
		CountHumanDNA:  10,
		Ratio:          0.5,
	}, nil
}
func (c *MockCollection) SaveHuman(*models.HumanModel) error { return nil }
func (c *MockCollection) Database() MongoDatabase            { return &MockDatabase{} }

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
