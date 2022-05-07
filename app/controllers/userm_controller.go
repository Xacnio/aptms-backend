package controllers

import (
	"database/sql"
	"fmt"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func MGetUser(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	userId, _ := strconv.ParseInt(c.Params("user_id"), 10, 32)
	result, err := db.MGetUser(int(buildingId), int(userId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MUserListType1(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	results, err := db.MGetUsersType1(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: results})
	}
}

func MUserListType2(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	results, err := db.MGetUsersType2(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: results})
	}
}

func MUserAddA(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
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
			Password:      utils.HashPassword(data.Email, "123"),
			Type:          0,
			CreatedBy:     user.ID,
			Email:         data.Email,
			PhoneNumber:   data.PhoneNumber,
			CreatedByName: fmt.Sprintf("%s %s", user.Name, user.Surname),
		}
		err := db.MCreateUserA(&userA)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		}
		db.MCreateUserAptPerm(userA.ID, uint(buildingId))
		return utils.ReturnSuccess(c, &models.ApiData{Result: userA})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func MUserKick(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseUint(c.Params("id"), 10, 32)
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
			ID uint `json:"id,string" validate:"required,numeric"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		exists, err := db.MUserCheckFlat(uint(buildingId), data.ID)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		}
		if exists {
			return utils.ReturnError(c, &models.ApiData{Error: "Error! The user is owner or tenant from a flat."})
		}
		err = db.MKickUser(uint(buildingId), data.ID)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}
