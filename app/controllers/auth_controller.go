package controllers

import (
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func LoginHandler(c *fiber.Ctx) error {
	payload := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: err, StatusCode: 400})
	}
	validate := validator.New()
	if err := validate.Struct(payload); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Invalid email or password", StatusCode: 400})
	}
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}

	user, err := db.GetUserByEmailPass(payload.Email, utils.HashPassword(payload.Email, payload.Password))
	if err != nil {
		if strings.Contains(err.Error(), "no rows in") {
			return utils.ReturnError(c, &models.ApiData{Error: "Invalid email or password", StatusCode: 200})
		}
		return utils.ReturnError(c, &models.ApiData{Error: "Database error" + err.Error(), StatusCode: 500})
	}
	if user.ID == 0 {
		return utils.ReturnError(c, &models.ApiData{Error: "Invalid email or password", StatusCode: 200})
	}
	tokenData, err := utils.CreateAccessToken(user.ID)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Access details could not be created", StatusCode: fiber.StatusInternalServerError})
	}
	ck := fiber.Map{
		"AccessToken": tokenData.AccessToken,
		"AccessUuid":  tokenData.AccessUuid,
		"UserID":      user.ID,
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: fiber.Map{"user": user, "tokens": ck}})
}

func LogoutHandler(c *fiber.Ctx) error {
	// TODO
	return utils.ReturnSuccess(c, &models.ApiData{})
}
