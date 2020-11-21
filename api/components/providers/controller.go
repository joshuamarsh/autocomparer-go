package providers

import (
	"carcompare/structs"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// GetAdverts get adverts from all sources provided
func (j *JSON) GetAdverts(c *fiber.Ctx) error {
	parameters := new(GetAdvertsParameters)
	if err := c.QueryParser(parameters); err != nil {
		j.logger.Errorf("%v", err)
	}

	parameters.Postcode = strings.ReplaceAll(parameters.Postcode, " ", "")

	adverts, err := j.providersManager.GetAdvert(parameters.Provider, parameters.Brand, parameters.Model, parameters.Postcode, parameters.Radius, parameters.SortBy)

	// j.logger.Debugf("%v", adverts)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err})
	}
	return c.JSON(adverts)
}

// GetCategories get adverts from all sources provided
func (j *JSON) GetMakes(c *fiber.Ctx) error {
	parameters := new(GetMakesParameters)
	if err := c.QueryParser(parameters); err != nil {
		j.logger.Errorf("%v", err)
	}

	categories, err := j.providersManager.GetMakes(parameters.Provider)

	// j.logger.Debugf("%v", adverts)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get makes", "data": err})
	}
	return c.JSON(categories)
}

func (j *JSON) GetModels(c *fiber.Ctx) error {
	var err error

	parameters := new(GetModelsParameters)
	if err := c.QueryParser(parameters); err != nil {
		j.logger.Errorf("%v", err)
	}

	models := []structs.ModelProvider{}
	if parameters.Brand != "" {
		models, err = j.providersManager.GetModels(parameters.Provider, parameters.Brand)
	}

	// j.logger.Debugf("%v", adverts)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err})
	}
	return c.JSON(models)
}
