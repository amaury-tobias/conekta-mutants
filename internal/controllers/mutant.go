package controllers

import (
	"github.com/amaury-tobias/conekta-mutants/internal/database"
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

	err = database.DBClient.Database("conekta").Collection("humans").SaveHuman(human)
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
	res, err := database.DBClient.Database("conekta").Collection("humans").GetStats()
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return c.JSON(res)
}
