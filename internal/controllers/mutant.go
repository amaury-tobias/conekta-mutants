package controllers

import (
	"github.com/amaury-tobias/conekta-mutants/internal/models"
	"github.com/amaury-tobias/conekta-mutants/internal/services"
	"github.com/gofiber/fiber/v2"
)

type MutantsController interface {
	IsMutant(c *fiber.Ctx) error
	GetStats(c *fiber.Ctx) error
}
type mutantsController struct {
	service services.MutantsService
}

func NewMutantsController(s services.MutantsService) MutantsController {
	return &mutantsController{service: s}
}

func (m *mutantsController) IsMutant(c *fiber.Ctx) error {
	human := new(models.HumanModel)
	if err := c.BodyParser(human); err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	isMutant, err := human.IsMutant()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: err.Error(),
		}
	}

	err = m.service.SaveHuman(human)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	if isMutant {
		return c.SendStatus(fiber.StatusOK)
	} else {
		return c.SendStatus(fiber.StatusForbidden)
	}
}
func (m *mutantsController) GetStats(c *fiber.Ctx) error {
	stats, err := m.service.GetStats()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}
	return c.JSON(stats)
}
