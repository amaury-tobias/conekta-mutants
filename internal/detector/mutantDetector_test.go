package detector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseDNA(t *testing.T) {
	type args struct {
		dna []string
	}
	type Test struct {
		name           string
		args           args
		want           [][]string
		wantErr        bool
		wantErrMessage string
	}

	tests := []Test{
		{
			name: "ParseDNA() valid",
			args: args{
				dna: []string{
					"ATGC",
					"TGGA",
					"GACC",
					"GCCT",
				},
			},
			want: [][]string{
				{"A", "T", "G", "C"},
				{"T", "G", "G", "A"},
				{"G", "A", "C", "C"},
				{"G", "C", "C", "T"},
			},
		},
		{
			name: "ParseDNA() invalid: no human",
			args: args{
				dna: []string{
					"ATGC",
					"TGGA",
					// "M" is not Valid
					"GACM",
					"GCCT",
				},
			},
			wantErr:        true,
			wantErrMessage: "la secuencia de ADN contiene bases no validas",
		},
		{
			name: "ParseDNA() invalid: bad sequencies size",
			args: args{
				dna: []string{
					"ATGC",
					"TGGA",
					"GCCT",
				},
			},
			wantErr:        true,
			wantErrMessage: "largo de secuencia invalido",
		},
	}

	for _, tt := range tests {
		got, err := ParseDNA(tt.args.dna)
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.EqualErrorf(t, err, tt.wantErrMessage, tt.name)
		}
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestGetColumns(t *testing.T) {
	type args struct {
		dnaSeq [][]string
	}
	type Test struct {
		name string
		args args
		want []string
	}
	tests := []Test{
		{
			name: "GetColumns()",
			args: args{
				dnaSeq: [][]string{
					{"A", "T", "G", "C"},
					{"T", "G", "G", "A"},
					{"G", "A", "C", "C"},
					{"G", "C", "C", "T"},
				},
			},
			want: []string{
				"ATGG",
				"TGAC",
				"GGCC",
				"CACT",
			},
		},
	}

	for _, tt := range tests {
		got := GetColumns(tt.args.dnaSeq)
		assert.NotNilf(t, got, tt.name)
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestSequenceIsMutant(t *testing.T) {
	type args struct {
		dna []string
	}
	type Test struct {
		name string
		args args
		want bool
	}

	tests := []Test{
		{
			name: "sequenceIsMutant() valid",
			args: args{
				dna: []string{"TTTT"},
			},
			want: true,
		},
		{
			name: "sequenceIsMutant() invalid",
			args: args{
				dna: []string{"AAASADAFA"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		got := SequenceIsMutant(tt.args.dna)
		assert.NotNilf(t, got, tt.name)
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestGetDiagonals(t *testing.T) {
	type args struct {
		dnaSeq [][]string
	}
	type Test struct {
		name string
		args args
		want []string
	}

	tests := []Test{
		{
			name: "GetDiagonals() valid",
			args: args{
				dnaSeq: [][]string{
					{"G", "T", "G", "C", "A"},
					{"G", "T", "G", "C", "A"},
					{"G", "T", "G", "C", "A"},
					{"G", "T", "G", "C", "A"},
					{"G", "T", "G", "C", "A"},
				},
			},
			// Only len(string) >= 4 are valid diagonals
			want: []string{
				"GTGC",
				"GTGCA",
				"TGCA",
				"TGCA",
				"GTGCA",
				"GTGC",
			},
		},
	}
	for _, tt := range tests {
		got := GetDiagonals(tt.args.dnaSeq)
		assert.NotNilf(t, got, tt.name)
		assert.Equalf(t, tt.want, got, tt.name)
	}
}

func TestIsMutant(t *testing.T) {
	type args struct {
		dna []string
	}
	type Test struct {
		name           string
		args           args
		want           bool
		wantErr        bool
		wantErrMessage string
	}

	tests := []Test{
		{
			name: "IsMutant() valid",
			args: args{
				dna: []string{
					"ATGC",
					"TAGA",
					"GAAC",
					"GCCA",
				},
			},
			want: true,
		},
		{
			name: "IsMutant() invalid: no human",
			args: args{
				dna: []string{
					"ATGC",
					"TGGA",
					// "M" is not Valid
					"GACM",
					"GCCT",
				},
			},
			wantErr:        true,
			wantErrMessage: "la secuencia de ADN contiene bases no validas",
		},
		{
			name: "IsMutant() invalid: bad sequencies size",
			args: args{
				dna: []string{
					"ATGC",
					"TAGA",
					"GCCA",
				},
			},
			wantErr:        true,
			wantErrMessage: "largo de secuencia invalido",
		},
	}

	for _, tt := range tests {
		got, err := IsMutant(tt.args.dna)
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)
		if tt.wantErr {
			assert.EqualErrorf(t, err, tt.wantErrMessage, tt.name)
		}
		assert.Equalf(t, tt.want, got, tt.name)
	}
}
