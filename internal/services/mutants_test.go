package services

import (
	"testing"

	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"github.com/amaury-tobias/conekta-mutants/internal/models"
	"github.com/stretchr/testify/assert"
)

func Test_mutantsService(t *testing.T) {
	type args struct {
		h *models.HumanModel
	}
	type Test struct {
		name    string
		m       MutantsService
		args    args
		wantErr bool
	}

	session := database.NewMockSession()
	assert.NotNilf(t, session, "Database Session")

	mutantsService, err := ServiceFromSession(session)
	assert.Nilf(t, err, "Mutants Service Error")
	assert.NotNilf(t, mutantsService, "Mutants Service Service Nil")

	tests := []Test{
		{
			name:    "mutantsService",
			m:       mutantsService,
			wantErr: false,
			args: args{
				h: &models.HumanModel{
					DNA: []string{
						"ATGC",
						"TGGA",
						"GACC",
						"GCCT",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		err := tt.m.SaveHuman(tt.args.h)
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)

		stats, err := tt.m.GetStats()
		assert.Equalf(t, tt.wantErr, err != nil, tt.name)
		assert.Equalf(t, &models.Stats{
			CountMutantDNA: 0,
			CountHumanDNA:  0,
			Ratio:          0,
		}, stats, tt.name)

	}
}
