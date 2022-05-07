package controllers

import (
	"database/sql"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/Xacnio/aptms-backend/pkg/utils"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func BuildingList(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	user := CheckUser(db, c)
	if user == nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: user not found", StatusCode: fiber.StatusInternalServerError})
	}
	if user.Type >= 1 {
		if len(c.Query("id")) > 0 {
			id, _ := strconv.ParseUint(c.Query("id"), 10, 16)
			if id > 0 {
				result, err := db.GetBuildingById(uint16(id))
				if err != nil && err != sql.ErrNoRows {
					return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
				} else {
					return utils.ReturnSuccess(c, &models.ApiData{Result: result})
				}
			}
		} else {
			results, err := db.GetBuildings(0, 10)
			if err != nil && err != sql.ErrNoRows {
				return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
			} else {
				return utils.ReturnSuccess(c, &models.ApiData{Result: results})
			}
		}
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func UpdateBuilding(c *fiber.Ctx) error {
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
			ID          int    `json:"id,string" validate:"required,numeric"`
			Name        string `json:"name" validate:"required"`
			CityID      uint8  `json:"city,string" validate:"required,numeric,gte=0"`
			DistrictID  uint16 `json:"district,string" validate:"required,numeric,gte=0"`
			Address     string `json:"address" validate:"required"`
			PhoneNumber string `json:"phoneNumber" validate:"omitempty,e164"`
			TaxNumber   string `json:"taxNumber" validate:""`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		var bui models.Building
		bui.ID = data.ID
		bui.Name = data.Name
		bui.City.ID = data.CityID
		bui.District.ID = data.DistrictID
		bui.Address = data.Address
		bui.PhoneNumber = data.PhoneNumber
		bui.TaxNumber = data.TaxNumber
		err := db.UpdateBuilding(bui)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func Blocks(c *fiber.Ctx) error {
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
			BuildingId int `json:"building_id,string" validate:"required,numeric"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		blocks, err := db.GetBlocks(data.BuildingId)
		if err != nil && err != sql.ErrNoRows {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		} else {
			return utils.ReturnSuccess(c, &models.ApiData{Result: blocks})
		}
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func AddBlock(c *fiber.Ctx) error {
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
			BuildingId int    `json:"building_id,string" validate:"required,numeric"`
			Letter     string `json:"letter" validate:"required,min=1,max=3"`
			DNumber    string `json:"d_number" validate:"required,min=1,max=3"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		block := models.Block{}
		block.BuildingID = data.BuildingId
		block.Letter = data.Letter
		block.DNumber = data.DNumber
		err := db.NewBuildingBlock(block)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Insert: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func AddBuilding(c *fiber.Ctx) error {
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
			CityID      uint8  `json:"city,string" validate:"required,numeric,gte=0"`
			DistrictID  uint16 `json:"district,string" validate:"required,numeric,gte=0"`
			Address     string `json:"address" validate:"required"`
			PhoneNumber string `json:"phoneNumber" validate:"omitempty,e164"`
			TaxNumber   string `json:"taxNumber" validate:""`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		var building models.Building
		building.Name = data.Name
		building.City.ID = data.CityID
		building.District.ID = data.DistrictID
		building.Address = data.Address
		building.PhoneNumber = data.PhoneNumber
		building.TaxNumber = data.TaxNumber
		err := db.NewBuilding(building)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Insert: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func DeleteBuilding(c *fiber.Ctx) error {
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
			BuildingId int `json:"building_id,string" validate:"required,numeric"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		err := db.DeleteBuilding(data.BuildingId)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		} else {
			return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
		}
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func DeleteBuildingBlock(c *fiber.Ctx) error {
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
			BuildingId int    `json:"building_id,string" validate:"required,numeric"`
			Letter     string `json:"letter" validate:"required,min=1,max=3"`
		})
		if err := c.BodyParser(&data); err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		validate := validator.New()
		errv := validate.Struct(data)
		if errv != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		err := db.DeleteBuildingBlock(data.BuildingId, data.Letter)
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
		} else {
			return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
		}
	}
	return utils.ReturnError(c, &models.ApiData{StatusCode: fiber.StatusUnauthorized})
}

func SystemInfo(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: db.SystemInfo()})
}
