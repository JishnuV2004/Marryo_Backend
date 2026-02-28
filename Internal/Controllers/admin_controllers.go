package controller

import (
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	services "marryo/Internal/Services"
	utils "marryo/Internal/Utils"
	"math"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AdminController struct {
	Service *services.AdminService
}

func NewAdminController(s *services.AdminService) *AdminController {
	return &AdminController{Service: s}
}

func (s *AdminController) GetAllProfiles(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	profiles,total, err := s.Service.GetAllProfiles(userID, page, limit)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(200).JSON(fiber.Map{
		"data": profiles,
		"pagination": fiber.Map{
			"page":        page,
			"limit":       limit,
			"total":       total,
			"total_pages": totalPages,
		},
	})
}

// Block Unblock
func (s *AdminController) BlockUnblock(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	status, err := s.Service.BlockUnblock(uint(id))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"status": status,
	})
}

// Count Users
func (s *AdminController) GetActiveUsersCount(c *fiber.Ctx) error {

	count, err := s.Service.GetActiveUsersCount()
	if err != nil {
		errormsg := utils.ErrorMessage(400, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errormsg,
		})
	}

	success := utils.SuccessResponse(count)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": success,
	})
}

// Count Matches
func (s *AdminController) GetUserMatchesCount(c *fiber.Ctx) error {

	count, err := s.Service.GetUserMatchesCount()
	if err != nil {
		errormsg := utils.ErrorMessage(400, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errormsg,
		})
	}

	success := utils.SuccessResponse(count)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": success,
	})
}

// DeleteProfileByID
func (s *AdminController) DeleteProfile(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	err = s.Service.DeleteProfile(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "deleted successfully",
	})
}

// GetMatches
func (s *AdminController) GetAllMatches(c *fiber.Ctx) error {

	matches, err := s.Service.GetAllMatches()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "failed to fetch matches",
		})
	}

	return c.JSON(fiber.Map{
		"count":   len(matches),
		"matches": matches,
	})
}

// Create User
func (s *AdminController) CreateUser(c *fiber.Ctx) error {

	var req models.User
	if err := c.BodyParser(&req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err := s.Service.CreateUser(&req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}
	success := utils.SuccessResponse("success")
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": success,
	})
}

//Edit Profile
func (s *AdminController) EditProfile(c *fiber.Ctx) error {

	idStr := c.Params("id")

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid user id",
		})
	}

	var input dto.EditProfile
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	editedprofile, err := s.Service.EditProfile(uint(id), &input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"err": err.Error(),
		})
	}

	success := utils.SuccessResponse(editedprofile)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": success,
	})
}

//Edit Admin Profile
func (s *AdminController) EditAdminProfile(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req dto.AdminProfile
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":"invalid input",
		})
	}

	if err := s.Service.EditAdminProfile(userID, &req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	success := utils.SuccessResponse("updated successfully")

	return  c.Status(200).JSON(fiber.Map{
		"message" : success,
	})
}	

//Change Admin Password
func (s *AdminController) ChangeAdminPassword(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	var req dto.AdminEditPassword
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error":"invalid input",
		})
	}

	if err := s.Service.ChangeAdminPassword(userID, req); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	success := utils.SuccessResponse("updated successfully")

	return  c.Status(200).JSON(fiber.Map{
		"message" : success,
	})
}

// ForgotPassword
func (s *AdminController) ForgotPassword(c *fiber.Ctx) error {

	userID := c.Locals("userID").(uint)

	if err := s.Service.ForgotPassword(userID); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "OTP sent successfully"})
}
