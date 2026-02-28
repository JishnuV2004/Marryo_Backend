package controller

import (
	// "fmt"
	"log"
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	"strconv"

	// models "marryo/Internal/Models"
	services "marryo/Internal/Services"
	utils "marryo/Internal/Utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	Services *services.UserService
}

func NewUserController(s *services.UserService) *UserController {
	return &UserController{Services: s}
}

// Profile
func (s *UserController) Profile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	user, err := s.Services.Profile(userID)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"err":   "profile not found",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		// "username": user.Username,
		"data": user,
	})
}

// EditProfile
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

// FilterProfiles
func (s *UserController) FilterProfiles(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req dto.SearchFilterRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	profiles, err := s.Services.FilterProfiles(&req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
			"err":   "filtering faild",
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": profiles,
	})
}

// FilterProfiles For User Home Page
func (s *UserController) HomeProfiles(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	profiles, err := s.Services.HomeProfiles(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "filtering failed",
			"err":   err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message":  "filtering successful",
		"profiles": profiles,
	})

}

// DeleteProfile
func (s *UserController) DeleteProfile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	err := s.Services.DeleteProfile(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "deleted successfully",
	})
}

// GetUserByID
func (s *UserController) GetProfile(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	profile, err := s.Services.GetProfile(uint(id))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"profile": profile,
	})
}

// SearchProfile
func (s *UserController) SearchProfiles(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req dto.SearchRequest
	if err := c.QueryParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid search query",
		})
	}

	profiles, err := s.Services.SearchProfiles(&req, userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": profiles,
	})
}

// ReceivedInterests
func (s *UserController) ReceivedInterests(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	interests, err := s.Services.ReceivedInterests(userID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"interests": interests,
	})
}

// Send Interest
func (s *UserController) SendInterest(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	interestID, _ := strconv.Atoi(c.Params("id"))

	// var req struct {
	// 	ReceiverID uint `json:"receiver_id"`
	// }

	// if err := c.BodyParser(&req); err != nil {
	// 	return c.Status(400).JSON(fiber.Map{
	// 		"error" : "invalid input",
	// 	})
	// }

	if err := s.Services.SendInterest(userID, uint(interestID)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "interest sented",
	})
}

// Accept Interest
func (s *UserController) AcceptInterest(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	interestID, _ := strconv.Atoi(c.Params("id"))

	if err := s.Services.AcceptInterest(userID, uint(interestID)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "interest accepted",
	})
}

// Decline Interest
func (s *UserController) DeclineInterest(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	interestID, _ := strconv.Atoi(c.Params("id"))

	if err := s.Services.DeclineInterest(userID, uint(interestID)); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "interest declined",
	})
}

// GetSentInterests
func (c *UserController) GetSentInterests(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	data, err := c.Services.GetInterests(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": data})
}

// ReceivedInterests
func (c *UserController) GetReceivedInterests(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	data, err := c.Services.GetReceivedInterests(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": data})
}

// GetAcceptedInterests
func (c *UserController) GetAcceptedInterests(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)

	data, err := c.Services.GetAcceptedInterests(userID)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{"data": data})
}

// internal/controllers/user_controller.go
func (c *UserController) SaveDeviceToken(ctx *fiber.Ctx) error {
	userID := ctx.Locals("userID").(uint)
	log.Printf("Received SaveDeviceToken request for userID: %d\n", userID)

	var body struct {
		Token    string `json:"token"`
		Platform string `json:"platform"`
	}

	if err := ctx.BodyParser(&body); err != nil {
		log.Printf("Error parsing SaveDeviceToken body for user %d: %v\n", userID, err)
		return ctx.Status(400).JSON(fiber.Map{"error": "invalid input"})
	}

	log.Printf("User %d registering FCM token: %s (%s)\n", userID, body.Token, body.Platform)

	token := models.DeviceToken{
		UserID:   userID,
		Token:    body.Token,
		Platform: body.Platform,
	}

	if err := c.Services.SaveDeviceToken(token); err != nil {
		log.Printf("Error saving token to DB for user %d: %v\n", userID, err)
		return ctx.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	log.Printf("Successfully saved device token for user %d\n", userID)
	return ctx.JSON(fiber.Map{"message": "device token saved"})
}

//UploadPhoto
// func (s *UserController) UploadPhoto(c *fiber.Ctx) error {

// 	userID := c.Locals("userID").(uint)

// 	imgURL := c.FormValue("image_url")
// 	if imgURL == "" {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error" : "image url required",
// 		})
// 	}

// 	if err := s.Services.UploadIMG(userID, imgURL); err != nil {
// 		return c.Status(400).JSON(fiber.Map{
// 			"error" : err.Error(),
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"message" : "image adedd",
// 	})
// }
