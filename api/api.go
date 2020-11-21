package api

import (
	"carcompare/api/components/providers"
	"carcompare/api/components/user"

	"github.com/apsdehal/go-logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// Run initialises endpoints and starts rest api
func Run(m providers.Manager, l *logger.Logger) {
	app := fiber.New()
	app.Use(cors.New())
	api := app.Group("/api")

	api.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Api Test")
	})

	providers.SetupRoutes(api, m, l)
	user.SetupRoutes(api)
	app.Listen(":8090")
}
