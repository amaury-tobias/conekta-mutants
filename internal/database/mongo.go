package database

import (
	"context"
	"strings"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Collection struct {
	collection *mongo.Collection
}

func (c *Collection) GetStats() (*models.Stats, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	groupTotalStage := bson.D{
		primitive.E{
			Key: "$group",
			Value: bson.D{
				primitive.E{
					Key:   "_id",
					Value: nil,
				},
				primitive.E{
					Key: "count",
					Value: bson.D{
						primitive.E{
							Key:   "$sum",
							Value: 1,
						},
					},
				},
				primitive.E{
					Key: "data",
					Value: bson.D{
						primitive.E{
							Key:   "$push",
							Value: "$$ROOT",
						},
					},
				},
			},
		},
	}
	unwindStage := bson.D{primitive.E{Key: "$unwind", Value: "$data"}}
	matchStage := bson.D{
		primitive.E{
			Key: "$match",
			Value: bson.D{
				primitive.E{
					Key:   "data.is_mutant",
					Value: true,
				},
			},
		},
	}
	groupStage := bson.D{
		primitive.E{
			Key: "$group",
			Value: bson.D{
				primitive.E{Key: "_id", Value: "$is_mutant"},
				primitive.E{
					Key:   "count_mutant_dna",
					Value: bson.D{primitive.E{Key: "$sum", Value: 1}},
				},
				primitive.E{
					Key:   "count_human_dna",
					Value: bson.D{primitive.E{Key: "$first", Value: "$count"}},
				},
			},
		},
	}
	projectStage := bson.D{
		primitive.E{
			Key: "$project",
			Value: bson.D{
				primitive.E{Key: "_id", Value: false},
				primitive.E{Key: "count_mutant_dna", Value: true},
				primitive.E{Key: "count_human_dna", Value: true},
				primitive.E{
					Key: "ratio",
					Value: bson.D{
						primitive.E{
							Key:   "$divide",
							Value: []string{"$count_mutant_dna", "$count_human_dna"},
						},
					},
				},
			},
		},
	}
	cursor, err := c.collection.Aggregate(ctx, mongo.Pipeline{
		groupTotalStage,
		unwindStage,
		matchStage,
		groupStage,
		projectStage,
	})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	result := new(models.Stats)
	for cursor.Next(ctx) {
		err := cursor.Decode(result)
		if err != nil {
			return nil, err
		}
	}
	if err := cursor.Err(); err != nil {
		return nil, err
	}
	return result, nil
}
func (c *Collection) SaveHuman(human *models.HumanModel) error {
	human.Key = strings.ToUpper(strings.Join(human.DNA, ""))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := c.collection.InsertOne(ctx, human)
	if err != nil {
		return err
	}
	return nil
}

func (c *Collection) Database() MongoDatabase {
	return &Database{
		database: c.collection.Database(),
	}
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
