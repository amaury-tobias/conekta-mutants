package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/amaury-tobias/conekta-mutants/internal/controllers"
)

func Init() *fiber.App {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
		ErrorHandler: func(c *fiber.Ctx, e error) error {
			code := fiber.StatusInternalServerError
			message := "Internal Server Error"

			if e, ok := e.(*fiber.Error); ok {
				code = e.Code
				message = e.Message
			}
			return c.Status(code).SendString(message)
		},
	})
	app.Use(cors.New())
	app.Use(recover.New())

	app.Post("/mutant", controllers.PostMutant)
	app.Get("/stats", controllers.GetStats)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).SendString("Not Found")
	})

	return app
}
