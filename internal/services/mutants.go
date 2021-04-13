package services

import (
	"context"
	"sync"
	"time"

	"github.com/amaury-tobias/conekta-mutants/internal/database"
	"github.com/amaury-tobias/conekta-mutants/internal/models"
	"github.com/amaury-tobias/conekta-mutants/internal/repository"
)

var once sync.Once

type MutantsService interface {
	SaveHuman(*models.HumanModel) error
	GetStats() (*models.Stats, error)
}

type mutantsService struct {
	repository repository.MutantsRepository
}

func (m *mutantsService) SaveHuman(h *models.HumanModel) error { return m.repository.SaveHuman(h) }
func (m *mutantsService) GetStats() (*models.Stats, error)     { return m.repository.GetStats() }

var instance *mutantsService

func NewMutantsService(r repository.MutantsRepository) MutantsService {
	once.Do(func() {
		instance = &mutantsService{
			repository: r,
		}
	})
	return instance
}

func ServiceFromSession(session database.Session) (MutantsService, error) {
	db, err := database.NewMongoDatabase(session)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer db.Client().Disconnect(ctx)
	mutantsRepository := repository.NewMutantsRepository(db)
	return NewMutantsService(mutantsRepository), nil
}
