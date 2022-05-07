package routes

import (
	"github.com/Xacnio/aptms-backend/app/controllers"
	"github.com/Xacnio/aptms-backend/pkg/configs"
	"github.com/Xacnio/aptms-backend/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	app.Use(cors.New(cors.Config{
		AllowOrigins: configs.Get("ALLOW_ORIGIN"),
	}))

	// PUBLIC REQUESTS
	api.Get("/info", controllers.SystemInfo)
	api.Get("/cities", controllers.GetCities)
	api.Get("/districts/:id", controllers.GetDistricts)

	// AUTHENTICATION
	auth := api.Group("/auth")
	auth.Post("/login", controllers.LoginHandler)
	auth.Get("/logout", controllers.LogoutHandler)

	// ADMIN REQUESTS
	admin := api.Group("/admin")
	admin.Use(middleware.Protected)

	// ADMIN USERS
	user := admin.Group("/user")
	user.Get("/", controllers.UserInfo)
	user.Get("/lista", controllers.UserListA)
	user.Post("/createa", controllers.UserAddA)
	user.Post("/updatea", controllers.UserEditA)

	// ADMIN BUILDINGS
	buildings := admin.Group("/buildings")
	buildings.Get("/", controllers.BuildingList)
	buildings.Post("/new", controllers.AddBuilding)
	buildings.Delete("/", controllers.DeleteBuilding)
	buildings.Post("/edit", controllers.UpdateBuilding)
	buildings.Post("/blocks", controllers.Blocks)
	buildings.Post("/blocks/new", controllers.AddBlock)
	buildings.Delete("/blocks", controllers.DeleteBuildingBlock)

	// MANAGER REQUESTS
	manager := api.Group("/manager")
	manager.Use(middleware.Protected)

	// BUILDING LIST FOR MANAGER
	mbuildings := manager.Group("/buildings")
	mbuildings.Get("/", controllers.MBuildingList)

	// BUILDING
	mbuilding := manager.Group("/buildings/:id")
	mbuilding.Use(middleware.ProtectedManager)
	mbuilding.Get("/details", controllers.MBuildingDetails)
	mbuilding.Get("/blocks", controllers.MGetBlocks)

	// BUILDING USERS
	mbuilding.Get("/users1", controllers.MUserListType1)
	mbuilding.Get("/users2", controllers.MUserListType2)
	mbuilding.Get("/user/:user_id", controllers.MGetUser)
	mbuilding.Post("/users/add", controllers.MUserAddA)
	mbuilding.Post("/users/kick", controllers.MUserKick)

	// BUILDING FLATS
	mbuilding.Get("/flats", controllers.MFlatList)
	mbuilding.Post("/flats", controllers.MAddFlat)
	mbuilding.Post("/flats/edit", controllers.MEditFlat)
	mbuilding.Delete("/flats", controllers.DeleteBuildingFlat)

	// BUILDING REVENUES
	mbuilding.Get("/revenues", controllers.MGetRevenues)
	mbuilding.Post("/revenues/new", controllers.MNewRevenue)
	mbuilding.Post("/revenues/edit/:rev_id", controllers.MEditRevenue)

	// BUILDING EXPENSES
	mbuilding.Get("/expenses", controllers.MGetExpenses)
	mbuilding.Post("/expenses/new", controllers.MNewExpense)
	mbuilding.Post("/expenses/edit/:exp_id", controllers.MEditExpense)
}
