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
	"time"
)

func MBuildingList(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	user := CheckUser(db, c)
	if user == nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	result, err := db.MGetBuildings(user.ID, 0, 30)
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MFlatList(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	result, err := db.MGetFlats(int(buildingId), 0, 30)
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MGetBlocks(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	result, err := db.MGetBlocks(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MAddFlat(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	data := new(struct {
		Type     uint8  `json:"type,string" validate:"omitempty,required,min=0,max=1"`
		BlockID  int    `json:"block_id,string" validate:"required,min=1"`
		Number   string `json:"number" validate:"required"`
		OwnerID  int    `json:"owner_id,string" validate:"omitempty,required,min=0"`
		TenantID int    `json:"tenant_id,string" validate:"omitempty,required,min=0"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	flat := models.Flat{}
	flat.BuildingID = int(buildingId)
	flat.BlockID = data.BlockID
	flat.OwnerID = data.OwnerID
	flat.Number = data.Number
	flat.TenantID = data.TenantID
	errt := db.MNewBuildingFlat(flat)
	if errt != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Insert: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
}

func DeleteBuildingFlat(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	data := new(struct {
		FlatId uint `json:"flat_id" validate:"required"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseUint(c.Params("id"), 10, 32)
	errt := db.MDeleteBuildingFlat(uint(buildingId), data.FlatId)
	if errt != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
}

func MEditFlat(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	}
	data := new(struct {
		ID       int `json:"id" validate:"required"`
		OwnerID  int `json:"owner_id,string" validate:"omitempty,required,min=0"`
		TenantID int `json:"tenant_id,string" validate:"omitempty,required,min=0"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	ownerIsMember, tenantIsMember := true, true
	if data.OwnerID > 0 {
		ov, err := db.IsBuildingMember(uint(data.OwnerID), uint(buildingId))
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		ownerIsMember = ov
	}
	if data.TenantID > 0 {
		tv, err := db.IsBuildingMember(uint(data.TenantID), uint(buildingId))
		if err != nil {
			return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
		}
		tenantIsMember = tv
	}
	if !tenantIsMember || !ownerIsMember {
		fmt.Println(tenantIsMember)
		fmt.Println(ownerIsMember)
		return utils.ReturnError(c, &models.ApiData{Error: "Sahip ya da kiracı bu apartmanın üyesi değil", StatusCode: fiber.StatusInternalServerError})
	}
	flat := models.Flat{}
	flat.BuildingID = int(buildingId)
	flat.ID = data.ID
	flat.OwnerID = data.OwnerID
	flat.TenantID = data.TenantID
	errt := db.MEditBuildingFlat(flat)
	if errt != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Edit: " + errt.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
}

func MGetRevenues(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	result, err := db.MRevenuesList(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MNewRevenue(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	data := new(struct {
		BlockID    int     `json:"block_id" validate:"required,numeric,min=1"`
		FlatID     int     `json:"flat_id" validate:"required,numeric,min=1"`
		PaidStatus bool    `json:"paid_status" validate:"omitempty,required"`
		PaidTime   string  `json:"paid_time" validate:"required"`
		PayerEmail string  `json:"payer_email" validate:"omitempty,required"`
		PayerName  string  `json:"payer_name" validate:"omitempty,required"`
		PayerPhone string  `json:"payer_phone" validate:"omitempty,required"`
		Total      float64 `json:"total" validate:"required"`
		Details    string  `json:"details" validate:"omitempty,required,max=128"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	rvn := new(models.Revenue)
	rvn.Payer.FullName = data.PayerName
	rvn.Payer.Email = data.PayerEmail
	rvn.Payer.Phone = data.PayerPhone
	rvn.Paid.Status = data.PaidStatus
	rvn.Details = data.Details
	if data.PaidTime != "NULL" {
		t, errs := time.Parse(time.RFC3339, data.PaidTime)
		if errs == nil {
			rvn.Paid.Time = &t
		}
	}
	rvn.Total = data.Total
	rvn.BuildingID = int(buildingId)
	rvn.FlatID.Valid = true
	rvn.FlatID.Int32 = int32(data.FlatID)
	err = db.MNewRevenue(rvn)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		if rvn.Paid.Status && rvn.Total > 0.0 {
			db.MBuildingCashIncrease(int(buildingId), rvn.Total)
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: rvn})
	}
}

func MEditRevenue(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	revId, _ := strconv.ParseInt(c.Params("rev_id"), 10, 32)
	data := new(struct {
		PaidStatus bool   `json:"paid_status" validate:"omitempty,required"`
		PaidTime   string `json:"paid_time" validate:"required"`
		PayerEmail string `json:"payer_email" validate:"omitempty,required"`
		PayerName  string `json:"payer_name" validate:"omitempty,required"`
		PayerPhone string `json:"payer_phone" validate:"omitempty,required"`
		Details    string `json:"details" validate:"omitempty,required,max=128"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	rvn := new(models.Revenue)
	rvn.ID = int(revId)
	rvn.Payer.FullName = data.PayerName
	rvn.Payer.Email = data.PayerEmail
	rvn.Payer.Phone = data.PayerPhone
	rvn.Paid.Status = data.PaidStatus
	rvn.Details = data.Details
	if data.PaidTime != "NULL" {
		t, errs := time.Parse(time.RFC3339, data.PaidTime)
		if errs == nil {
			rvn.Paid.Time = &t
		}
	}
	rvn.BuildingID = int(buildingId)
	err = db.MEditRevenue(rvn)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		if rvn.Paid.Status && rvn.Total > 0.0 {
			db.MBuildingCashIncrease(int(buildingId), rvn.Total)
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
}

func MBuildingDetails(c *fiber.Ctx) error {
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	result, err := db.MBuildingDetails(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MGetExpenses(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	result, err := db.MExpensesList(int(buildingId))
	if err != nil && err != sql.ErrNoRows {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		return utils.ReturnSuccess(c, &models.ApiData{Result: result})
	}
}

func MNewExpense(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	data := new(struct {
		PaidStatus bool    `json:"paid_status" validate:"omitempty,required"`
		PaidTime   string  `json:"paid_time" validate:"required"`
		ToEmail    string  `json:"to_email" validate:"omitempty,required"`
		ToName     string  `json:"to_name" validate:"required"`
		ToPhone    string  `json:"to_phone" validate:"omitempty,required"`
		Total      float64 `json:"total" validate:"required"`
		Details    string  `json:"details" validate:"omitempty,required,max=128"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	exp := new(models.Expense)
	exp.To.FullName = data.ToName
	exp.To.Email = data.ToEmail
	exp.To.Phone = data.ToPhone
	exp.Paid.Status = data.PaidStatus
	exp.Details = data.Details
	if data.PaidTime != "NULL" {
		t, errs := time.Parse(time.RFC3339, data.PaidTime)
		if errs == nil {
			exp.Paid.Time = &t
		}
	}
	exp.Total = data.Total
	exp.BuildingID = int(buildingId)
	err = db.MNewExpense(exp)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		if exp.Paid.Status && exp.Total > 0.0 {
			db.MBuildingCashDecrease(int(buildingId), exp.Total)
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: exp})
	}
}

func MEditExpense(c *fiber.Ctx) error {
	db, err := database.OpenDBConnection()
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	buildingId, _ := strconv.ParseInt(c.Params("id"), 10, 32)
	expId, _ := strconv.ParseInt(c.Params("exp_id"), 10, 32)
	data := new(struct {
		PaidStatus bool   `json:"paid_status" validate:"omitempty,required"`
		PaidTime   string `json:"paid_time" validate:"required"`
		ToEmail    string `json:"to_email" validate:"omitempty,required"`
		ToName     string `json:"to_name" validate:"omitempty,required"`
		ToPhone    string `json:"to_phone" validate:"omitempty,required"`
		Details    string `json:"details" validate:"omitempty,required,max=128"`
	})
	if err := c.BodyParser(&data); err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Parser: " + err.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	validate := validator.New()
	errv := validate.Struct(data)
	if errv != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Validate: " + errv.Error(), StatusCode: fiber.StatusInternalServerError})
	}
	exp := new(models.Expense)
	exp.ID = int(expId)
	exp.To.FullName = data.ToName
	exp.To.Email = data.ToEmail
	exp.To.Phone = data.ToPhone
	exp.Details = data.Details
	exp.Paid.Status = data.PaidStatus
	if data.PaidTime != "NULL" {
		t, errs := time.Parse(time.RFC3339, data.PaidTime)
		if errs == nil {
			exp.Paid.Time = &t
		}
	}
	exp.BuildingID = int(buildingId)
	err = db.MEditExpense(exp)
	if err != nil {
		return utils.ReturnError(c, &models.ApiData{Error: "Database error", StatusCode: fiber.StatusInternalServerError})
	} else {
		if exp.Paid.Status && exp.Total > 0.0 {
			db.MBuildingCashDecrease(int(buildingId), exp.Total)
		}
		return utils.ReturnSuccess(c, &models.ApiData{Result: "OK"})
	}
}
