package controller

import (
	// "fmt"
	dto "marryo/Internal/DTO"
	services "marryo/Internal/Services"
	utils "marryo/Internal/Utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Services *services.UserService
}

func NewUserController(s *services.UserService)  *UserController {
	return &UserController{Services: s}
}

//Profile
func (s *UserController) Profile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	user, err := s.Services.Profile(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"err": "profile not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		// "username": user.Username,
		"data": user,
	})
}

//EditProfile
func (s *UserController) EditProfile(c *fiber.Ctx) error {
	userID, ok := c.Locals("userID").(uint)
	if !ok {
		return c.SendStatus(http.StatusUnauthorized)
	}

	var input dto.EditProfile
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	updated, err := s.Services.EditProfile(userID, &input)
	if err != nil {
		return c.Status(400).JSON(utils.ErrorMessage(400, err))
	}

	return c.JSON(utils.SuccessResponseMsg(updated, "Updated successfully"))
}


//FilterProfiles
func (s *UserController) FilterProfiles(c *fiber.Ctx) error {

	var req dto.SearchFilterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	profiles, err := s.Services.FilterProfiles(&req)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"err" : "filtering faild",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": profiles,
	})
}

//FilterProfiles For User Home Page
func (s *UserController) HomeProfiles(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	profiles, err := s.Services.HomeProfiles(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error" : "filtering failed",
			"err" : err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "filtering successful",
		"profiles" : profiles,
	})

}

//DeleteProfile
func (s *UserController) DeleteProfile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	err := s.Services.DeleteProfile(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err" : err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message" : "deleted successfully",
	})
}