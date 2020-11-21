package providers

import (
	"encoding/json"
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cache"
)

// JSON logger and providers structure
type JSON struct {
	logger           *logger.Logger
	providersManager Manager
}

// SetupRoutes initialises adverts routes
func SetupRoutes(app fiber.Router, m Manager, l *logger.Logger) {
	j := JSON{
		providersManager: m,
		logger:           l,
	}

	advertRoute := app.Group("/adverts")
	advertRoute.Get("/", cache.New(cache.Config{
		Expiration: 30 * time.Minute,
		Key: func(c *fiber.Ctx) string {
			parameters := new(GetAdvertsParameters)
			if err := c.QueryParser(parameters); err != nil {
				j.logger.Errorf("%s", err)
				return err.Error()
			}

			parametersJSON, err := json.Marshal(parameters)
			if err != nil {
				j.logger.Errorf("%s", err)
				return err.Error()
			}

			return string(parametersJSON)
		},
	}), j.GetAdverts)

	app.Get("/GetMakes", cache.New(cache.Config{
		Expiration: 30 * time.Minute,
	}), j.GetMakes)

	app.Get("/GetModels", cache.New(cache.Config{
		Expiration: 30 * time.Minute,
		Key: func(c *fiber.Ctx) string {
			parameters := new(GetModelsParameters)
			if err := c.QueryParser(parameters); err != nil {
				j.logger.Errorf("%v", err)
			}
			return parameters.Brand
		},
	}), j.GetModels)

	// app.Get("/GetModels", cache.New(), j.GetModels)
}
