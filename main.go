package main

import (
	"github.com/Xacnio/aptms-backend/pkg/routes"
	"github.com/Xacnio/aptms-backend/platform/database"
	"github.com/gofiber/fiber/v2"
)

func main() {
	//fmt.Println(utils.HashPassword("test2@test.com", "123"))

	db, err := database.OpenDBConnection()
	if err != nil {
		panic(err)
	}

	go func() {
		db.CheckTablesExist()
		//db.TestData()
	}()

	app := fiber.New()

	routes.SetupRoutes(app)

	err = app.Listen(":80")
	if err != nil {
		panic(err)
	}
}
