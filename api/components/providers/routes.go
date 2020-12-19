package providers

import (
	"time"

	"github.com/apsdehal/go-logger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
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

	advertRoute.Get("/", j.GetAdverts)

	app.Get("/GetMakes", limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
	}), j.GetMakes)

	app.Get("/GetModels", limiter.New(limiter.Config{
		Max:        30,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusTooManyRequests)
		},
	}), j.GetModels)
}
