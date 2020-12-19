package providers

import (
	"carcompare/cache"
	"carcompare/config"
	"carcompare/structs"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

// GetAdverts get adverts from all providers provided
func (j *JSON) GetAdverts(c *fiber.Ctx) error {
	parameters := new(GetAdvertsParameters)
	if err := c.QueryParser(parameters); err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err.Error()})
	}

	parameters.Postcode = strings.ReplaceAll(parameters.Postcode, " ", "")
	cacheKey := structs.Sha256(parameters)

	if config.Config("CACHE") == "true" {
		advertsCacheString, err := cache.RedisDB.Get(ctx, cacheKey).Result()
		if err != nil {
			j.logger.Debugf("%v", err)
		}

		if advertsCacheString != "" {
			var advertsCache structs.AdvertProviders
			if err := json.Unmarshal([]byte(advertsCacheString), &advertsCache); err != nil {
				return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err.Error()})
			}

			return c.JSON(advertsCache)
		}
	}

	adverts, err := j.providersManager.GetAdvert(parameters.Provider, parameters.Brand, parameters.Model, parameters.Postcode, parameters.Radius, parameters.SortBy)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err.Error()})
	}

	if config.Config("CACHE") == "true" {
		advertsJSON, err := json.Marshal(&adverts)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err.Error()})
		}

		err = cache.RedisDB.SetEX(ctx, cacheKey, string(advertsJSON), 30*time.Minute).Err()
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get adverts", "data": err.Error()})
		}
	}
	return c.JSON(adverts)
}

// GetMakes get makes from all providers provided
func (j *JSON) GetMakes(c *fiber.Ctx) error {
	parameters := new(GetMakesParameters)
	if err := c.QueryParser(parameters); err != nil {
		j.logger.Errorf("%v", err)
	}

	makesCacheString, err := cache.RedisDB.Get(ctx, "makes").Result()
	if err != nil {
		j.logger.Debugf("%v", err)
	}

	if makesCacheString != "" {
		var makesCache []structs.MakeProvider
		if err := json.Unmarshal([]byte(makesCacheString), &makesCache); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err.Error()})
		}

		return c.JSON(makesCache)
	}

	makes, err := j.providersManager.GetMakes(parameters.Provider)

	// j.logger.Debugf("%v", adverts)
	makesJSON, err := json.Marshal(&makes)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get makes", "data": err.Error()})
	}

	err = cache.RedisDB.SetEX(ctx, "makes", string(makesJSON), 30*time.Minute).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get makes", "data": err.Error()})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get makes", "data": err})
	}
	return c.JSON(makes)
}

// GetModels get models from all providers provided
func (j *JSON) GetModels(c *fiber.Ctx) error {
	var err error

	parameters := new(GetModelsParameters)
	if err := c.QueryParser(parameters); err != nil {
		j.logger.Errorf("%v", err)
	}

	modelsCacheString, err := cache.RedisDB.Get(ctx, parameters.Brand).Result()
	if err != nil {
		j.logger.Debugf("%v", err)
	}

	if modelsCacheString != "" {
		var modelsCache []structs.ModelProvider
		if err := json.Unmarshal([]byte(modelsCacheString), &modelsCache); err != nil {
			return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err.Error()})
		}

		return c.JSON(modelsCache)
	}

	models := []structs.ModelProvider{}
	if parameters.Brand != "" {
		models, err = j.providersManager.GetModels(parameters.Provider, parameters.Brand)
	}

	modelsJSON, err := json.Marshal(&models)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err.Error()})
	}

	err = cache.RedisDB.SetEX(ctx, parameters.Brand, string(modelsJSON), 30*time.Minute).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err.Error()})
	}

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get models", "data": err})
	}
	return c.JSON(models)
}
