package user

import "github.com/gofiber/fiber/v2"

// SetupRoutes initialises user routes
func SetupRoutes(app fiber.Router) {
	userRoute := app.Group("/user")

	userRoute.Post("/register", CreateUser)

	userRoute.Get("/login", func(c *fiber.Ctx) error {
		return c.SendString("Login User")
	})
}
