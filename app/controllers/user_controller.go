package controllers

import (
	"database/sql"
	"fmt"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func UserInfo(c *fiber.Ctx) error {
	userid := c.Locals("User-ID")

	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}

	user, err := db.GetUserById(userid.(uint))
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error"})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: user})
}

func CheckUser(db *database.Queries, c *fiber.Ctx) *models.User {
	userid := c.Locals("User-ID")
	user, err := db.GetUserById(userid.(uint))
	if err != nil {
		return nil
	}
	return &user
}

func UserListA(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	user := CheckUser(db, c)
	if user == nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	if user.Type >= 1 {
		results, err := db.GetUsersA(0, 100)
		if err != nil && err != sql.ErrNoRows {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		} else {
			return utils.ReturnSuccess(c, &models.ApiData{Result: results})
		}
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func UserAddA(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	user := CheckUser(db, c)
	if user == nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	if user.Type >= 1 {
		data := new(struct {
			Name        string `json:"name" validate:"required"`
			Surname     string `json:"surname" validate:"required"`
			Email       string `json:"email" validate:"required,email"`
			PhoneNumber string `json:"phoneNumber" validate:"omitempty,required,e164"`
			Buildings   []uint `json:"buildings" validate:"dive,required,min=1"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		userA := models.User{
			Name:          data.Name,
			Surname:       data.Surname,
			Type:          0,
			CreatedBy:     user.ID,
			Email:         data.Email,
			PhoneNumber:   data.PhoneNumber,
			CreatedByName: fmt.Sprintf("%s %s", user.Name, user.Surname),
		}
		err := db.CreateUserA(&userA)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		}
		db.CreateUserBuildingPerm(userA.ID, data.Buildings, 1)
		return utils.ReturnSuccess(c, &models.ApiData{Result: userA})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func UserEditA(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	user := CheckUser(db, c)
	if user == nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	if user.Type >= 1 {
		data := new(struct {
			ID          uint   `json:"id,string" validate:"required,number"`
			Name        string `json:"name" validate:"required"`
			Surname     string `json:"surname" validate:"required"`
			Email       string `json:"email" validate:"required,email"`
			PhoneNumber string `json:"phoneNumber" validate:"omitempty,required,e164"`
			Buildings   []uint `json:"buildings" validate:"omitempty,dive"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		userA := models.User{
			ID:          data.ID,
			Name:        data.Name,
			Surname:     data.Surname,
			Email:       data.Email,
			PhoneNumber: data.PhoneNumber,
		}
		err := db.UpdateUserA(userA)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		db.UpdateUserAptPerm(userA.ID, data.Buildings, 1)
		return utils.ReturnSuccess(c, &models.ApiData{Result: userA})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}
