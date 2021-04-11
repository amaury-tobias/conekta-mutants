package detector

import (
	"errors"
	"regexp"
	"strings"
)

var patternParse *regexp.Regexp = regexp.MustCompile("(?i)[^ATCG]+")
var patternSequence *regexp.Regexp = regexp.MustCompile("(A{4,}|T{4,}|G{4,}|C{4,})")

func ParseDNA(dna []string) ([][]string, error) {
	dnaSeq := make([][]string, len(dna))

	for i, d := range dna {
		if len(d) != len(dna) {
			return nil, errors.New("largo de secuencia invalido")
		} else if patternParse.MatchString(d) {
			return nil, errors.New("la secuencia de ADN contiene bases no validas")
		}
		dnaSeq[i] = strings.Split(strings.ToUpper(d), "")
	}
	return dnaSeq, nil
}

func GetColumns(dnaSeq [][]string) []string {
	columns := make([]string, len(dnaSeq))
	for _, dna := range dnaSeq {
		for i, c := range dna {
			columns[i] += c
		}
	}
	return columns
}

func SequenceIsMutant(dna []string) bool {
	for _, sequence := range dna {
		if isMutant := patternSequence.MatchString(sequence); isMutant {
			return true
		}
	}
	return false
}

func GetDiagonals(dnaSeq [][]string) []string {
	diagonals := make([]string, 0)
	acc := ""

	for k := 0; k < (len(dnaSeq)*2)-1; k++ {
		for j := 0; j <= k; j++ {
			i := k - j
			if i < len(dnaSeq) && j < len(dnaSeq) {
				acc += dnaSeq[i][j]
			}
		}
		if len(acc) > 3 {
			diagonals = append(diagonals, acc)
		}
		acc = ""
	}

	for n := -len(dnaSeq); n <= len(dnaSeq); n++ {
		for i := 0; i < len(dnaSeq); i++ {
			if (i-n >= 0) && (i-n < len(dnaSeq)) {
				acc += dnaSeq[i][i-n]
			}
		}
		if len(acc) > 3 {
			diagonals = append(diagonals, acc)
		}
		acc = ""
	}

	return diagonals
}

func IsMutant(dna []string) (bool, error) {
	dnaSequence, err := ParseDNA(dna)
	if err != nil {
		return false, err
	}

	dnaSequences := append(
		// Default DNA used as rows
		dna,
		append(
			GetColumns(dnaSequence),
			GetDiagonals(dnaSequence)...,
		)...,
	)

	return SequenceIsMutant(dnaSequences), nil
}
