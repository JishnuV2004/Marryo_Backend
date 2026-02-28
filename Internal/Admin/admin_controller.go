package admin

import (
	"errors"
	repositories "marryo/Internal/Repositories"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
)

type AdminController struct{
	repo *repositories.Repository
}

func NewAdminController(r *repositories.Repository) *AdminController {
	return &AdminController{repo: r}
}

// Temporary hardcoded admin
var adminEmail = "admin@matrimony.com"
var adminPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("admin123"), 14)

func (r *AdminController) ShowLogin(c *fiber.Ctx) error {
	return c.Render("Admin/login", fiber.Map{})
}

func (r *AdminController) AdminLogin(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {

		email := c.FormValue("email")
		password := c.FormValue("password")

		if err := validateAdmin(email, password); err != nil {
			return c.Render("Admin/login", fiber.Map{
				"Error": "Invalid email or password",
			})
		}

		sess, _ := store.Get(c)
		sess.Set("admin_email", email)
		sess.Save()

		return c.Redirect("/admin/dashboard")
	}
}

func (r *AdminController) Dashboard(c *fiber.Ctx) error {

	email := c.Locals("admin_email")

	data := fiber.Map{
		"Email":        email,
		"PageTitle":    "Dashboard",
		"Active":       "dashboard",
		"TotalUsers":   1200,
		"TotalMatches": 340,
		"TotalReports": 12,
	}

	return c.Render("Admin/dashboard", data, "Layouts/admin_layout")
}


func (r *AdminController) Logout(store *session.Store) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sess, _ := store.Get(c)
		sess.Destroy()
		return c.Redirect("/Admin/login")
	}
}

// Internal validation function
func validateAdmin(email, password string) error {

	if email != adminEmail {
		return errors.New("invalid")
	}

	err := bcrypt.CompareHashAndPassword(adminPasswordHash, []byte(password))
	if err != nil {
		return errors.New("invalid")
	}

	return nil
}