package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func AdminRoutes(app *fiber.App, store *session.Store, adminController *AdminController) {

	admin := app.Group("/admin")

	admin.Get("/login", adminController.ShowLogin)
	admin.Post("/login", adminController.AdminLogin(store))
	admin.Get("/dashboard", adminController.Dashboard)
	admin.Get("/logout", adminController.Logout(store))

	// admin.Get("/dashboard",
	// 	adminMiddleware.AdminProtected(store),
	// 	adminController.Dashboard,
	// )

	// admin.Get("/logout",
	// 	adminMiddleware.AdminProtected(store),
	// 	adminController.Logout(store),
	// )

}