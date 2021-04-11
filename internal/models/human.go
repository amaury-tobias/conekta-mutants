package models

import (
	"context"
	"strings"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"github.com/amaury-tobias/conekta-mutants/internal/detector"
)

type HumanModel struct {
	Key    string   `json:"key" bson:"_id"`
	DNA    []string `json:"dna" bson:"dna"`
	Mutant bool     `json:"is_mutant" bson:"is_mutant"`
}

func (hm *HumanModel) Save() error {
	hm.Key = strings.ToUpper(strings.Join(hm.DNA, ""))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := database.MutantCollection.InsertOne(ctx, hm)
	if err != nil {
		return err
	}
	return nil
}

func (h *HumanModel) IsMutant() (bool, error) {
	got, err := detector.IsMutant(h.DNA)
	if err != nil {
		return false, err
	}
	h.Mutant = got
	return got, nil
}
