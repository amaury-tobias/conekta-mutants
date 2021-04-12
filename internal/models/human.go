package models

import (
	"github.com/amaury-tobias/conekta-mutants/internal/detector"
)

type HumanModel struct {
	Key    string   `json:"key" bson:"_id"`
	DNA    []string `json:"dna" bson:"dna"`
	Mutant bool     `json:"is_mutant" bson:"is_mutant"`
}

func (h *HumanModel) IsMutant() (bool, error) {
	got, err := detector.IsMutant(h.DNA)
	if err != nil {
		return false, err
	}
	h.Mutant = got
	return got, nil
}
