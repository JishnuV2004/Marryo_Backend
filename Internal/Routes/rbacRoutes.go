package routes

import (
	consts "marryo/Internal/Consts"
	controller "marryo/Internal/Controllers"
	middleware "marryo/Internal/MiddleWare"

	"github.com/gofiber/fiber/v2"
)

func RBACRoutes(app *fiber.App, rbac *controller.RoleHandler) {

	admin := app.Group("/admin", middleware.AuthMiddleWare())

	admin.Post("/createrole", middleware.RolesMiddleWare(consts.Admin), rbac.CreateRole)
	admin.Put("/update", middleware.RolesMiddleWare(consts.Admin), rbac.UpdateRolePermissions)
	admin.Get("/permissions/:id", middleware.RolesMiddleWare(consts.Admin), rbac.GetRolePermissions)
	admin.Get("/allpermissions", middleware.RolesMiddleWare(consts.Admin, consts.Staff), rbac.GetAllPermissions)

}
