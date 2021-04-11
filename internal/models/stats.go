package models

import (
	"context"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Stats struct {
	CountMutantDNA int     `json:"count_mutant_dna" bson:"count_mutant_dna"`
	CountHumanDNA  int     `json:"count_human_dna" bson:"count_human_dna"`
	Ratio          float64 `json:"ratio" bson:"ratio"`
}

func (s *Stats) GetStats() (*Stats, error) {
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

	cursor, err := database.MutantCollection.Aggregate(ctx, mongo.Pipeline{
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
	result := new(Stats)

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
