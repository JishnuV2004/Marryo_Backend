package routes

import (
	consts "marryo/Internal/Consts"
	controller "marryo/Internal/Controllers"
	middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app *fiber.App, userController *controller.UserController) {

	user := app.Group("/user")
	user.Use(middleware.RolesMiddleWare(consts.User))

	user.Get("/profile", userController.Profile)
	user.Post("/editprofile", userController.EditProfile)
	user.Post("/filterprofiles", userController.FilterProfiles)
	user.Get("/homeprofiles", userController.HomeProfiles)
	user.Post("/deleteprofile", userController.DeleteProfile)
	user.Get("/getprofile/:id", userController.GetProfile)
	user.Post("/searchprofiles", userController.SearchProfiles)
	user.Get("/receivedinterests", userController.ReceivedInterests)
	user.Post("/sendinterest/:id", userController.SendInterest)
	user.Post("/acceptinterest/:id", userController.AcceptInterest)
	user.Post("/declinedinterest/:id", userController.DeclineInterest)
	user.Get("/getinterests", userController.GetSentInterests)
	user.Get("/getreceived", userController.GetReceivedInterests)
	user.Get("/getaccepts", userController.GetAcceptedInterests)
	user.Post("/device-token", userController.SaveDeviceToken)
	// user.Post("/uploadphoto", userController.UploadPhoto)

}
