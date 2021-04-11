package controllers

import (
	"github.com/amaury-tobias/conekta-mutants/internal/models"
	"github.com/gofiber/fiber/v2"
)

func PostMutant(c *fiber.Ctx) error {
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

	err = human.Save()
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

func GetStats(c *fiber.Ctx) error {
	stats := new(models.Stats)
	res, err := stats.GetStats()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(res)
}
