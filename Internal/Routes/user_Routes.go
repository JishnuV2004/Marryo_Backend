package routes

import (
	controller "marryo/Internal/Controllers"
	middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userController *controller.UserController){

	user := app.Group("/user")
	user.Use(middleware.MiddleWare())

	user.Get("/profile", userController.Profile)
	user.Post("/editprofile", userController.EditProfile)
	user.Post("/filterprofiles", userController.FilterProfiles)
	user.Get("/homeprofiles", userController.HomeProfiles)
	user.Post("/deleteprofile", userController.DeleteProfile)
		
}