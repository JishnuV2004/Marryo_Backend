package routes

import (
	consts "marryo/Internal/Consts"
	controller "marryo/Internal/Controllers"
	middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
)

func AdminRoutes(app *fiber.App, adminController *controller.AdminController, authController *controller.AuthController, userController *controller.UserController) {

	admin := app.Group("/admin")

	// admin.Post("/login", authController.Login)

	admin.Use(middleware.AuthMiddleWare(), middleware.RolesMiddleWare(consts.Admin, consts.Staff))

	admin.Post("/logout", authController.Logout)
	admin.Get("/getallprofiles", adminController.GetAllProfiles)
	admin.Get("/getuser/:id", userController.GetProfile)
	admin.Post("/blockunblock/:id", adminController.BlockUnblock)
	admin.Get("/getprofile/:id", userController.GetProfile)
	admin.Get("/countusers", adminController.GetActiveUsersCount)
	admin.Get("/countmatches", adminController.GetUserMatchesCount)
	admin.Post("/deleteprofile/:id", adminController.DeleteProfile)
	admin.Get("/getallmatches", adminController.GetAllMatches)
	admin.Post("/adduser", adminController.CreateUser)
	admin.Post("/logout", authController.Logout)
	admin.Post("/editprofile/:id", adminController.EditProfile)
	admin.Post("/editadminprofile", adminController.EditAdminProfile)
	admin.Post("/changepassword", adminController.ChangeAdminPassword)
	admin.Post("/sentotp", adminController.ForgotPassword)
}
