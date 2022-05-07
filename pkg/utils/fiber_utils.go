package utils

import (
	"encoding/json"
	"errors"
	"github.com/Xacnio/aptms-backend/app/models"
	"github.com/gofiber/fiber/v2"
	"reflect"
	"strings"
)

func GetBearerToken(c *fiber.Ctx) string {
	bearToken := c.Get(fiber.HeaderAuthorization)
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func ReturnSuccess(c *fiber.Ctx, data *models.ApiData) error {
	data.IsSuccess = true
	if data.StatusCode == 0 {
		data.StatusCode = 200
	}
	if data.Headers == nil {
		data.Headers = []map[string]string{}
	}

	jsonValue, _ := json.Marshal(data)
	return c.Status(data.StatusCode).Send(jsonValue)
}

func ReturnError(c *fiber.Ctx, data *models.ApiData) error {
	data.IsSuccess = false
	if data.StatusCode == 0 {
		data.StatusCode = 200
	}
	if data.Headers == nil {
		data.Headers = []map[string]string{}
	}
	if reflect.TypeOf(data.Error) == reflect.TypeOf(errors.New("")) {
		data.Error = data.Error.(error).Error()
	}
	if reflect.TypeOf(data.Error) != reflect.TypeOf("") {
		data.Error = nil
	}

	jsonValue, _ := json.Marshal(data)
	return c.Status(data.StatusCode).Send(jsonValue)
}
