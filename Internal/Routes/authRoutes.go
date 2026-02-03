package routes

import (
	controller "marryo/Internal/Controllers"
	// middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
)


func Routes(app *fiber.App, controller *controller.AuthController){

	app.Post("/signup", controller.Signup)
	app.Post("/verifyotp", controller.VerifyOTP)
	app.Post("/completesignup", controller.CompleteSignup)
	app.Post("/login", controller.Login)
	app.Post("/logout", controller.Logout)
	app.Post("/refresh", controller.Refresh)
	// app.Get("/profile", middleware.MiddleWare(), controller.Profile)

}