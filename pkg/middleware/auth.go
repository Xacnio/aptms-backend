package middleware

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

// Protected protect routes
func Protected(c *fiber.Ctx) error {
	bearerToken := utils.GetBearerToken(c)
	if len(bearerToken) < 1 {
		return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized, Error: "Not logged in (1)"})
	}
	token, err := utils.VerifyAccessToken(bearerToken)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized, Error: "Not logged in (2)"})
	}
	authInfo, err := utils.ExtractTokenPayload(token)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized, Error: "Not logged in (3)"})
	}
	c.Locals("User-ID", authInfo.UserId)
	return c.Next()
}
