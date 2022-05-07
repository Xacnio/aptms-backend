package middleware

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func ProtectedManager(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
	}
	userId := c.Locals("User-ID")
	access, _ := db.IsBuildingManager(userId.(uint), uint(buildingId))
	if access == true {
		return c.Next()
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}
