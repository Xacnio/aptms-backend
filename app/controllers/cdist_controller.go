package controllers

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetCities(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}

	cities, err := db.GetCities()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: cities})
}

func GetDistricts(c *fiber.Ctx) error {
	cityId := c.Params("id")

	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}
	cityIdQ, _ := strconv.ParseUint(cityId, 10, 8)
	districts, err := db.GetDistricts(uint8(cityIdQ))
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: districts})
}
