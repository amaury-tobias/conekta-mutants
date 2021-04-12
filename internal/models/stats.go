package models

type Stats struct {
	CountMutantDNA int     `json:"count_mutant_dna" bson:"count_mutant_dna"`
	CountHumanDNA  int     `json:"count_human_dna" bson:"count_human_dna"`
	Ratio          float64 `json:"ratio" bson:"ratio"`
}
