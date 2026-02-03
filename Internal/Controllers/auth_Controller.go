package controller

import (
	dto "marryo/Internal/DTO"
	// models "marryo/Internal/Models"
	services "marryo/Internal/Services"
	utils "marryo/Internal/Utils"
	"net/http"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v5"
)

type AuthController struct {
	Service *services.AuthService
}

func NewAuthController(s *services.AuthService) *AuthController {
	return &AuthController{Service: s}
}

// Signup Func
func (s *AuthController) Signup(c *fiber.Ctx) error {
	var newUser dto.RegisterRequest
	if err := c.BodyParser(&newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	if err := utils.Validator.Struct(newUser); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	data, err := s.Service.Signup(&newUser)
	if err != nil {
		errormsg := utils.ErrorMessage(400, err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": errormsg,
		})
	}

	success := utils.SuccessResponse(data)
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message":success,
	})
}

//VerifyOTP
func (s *AuthController) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	if err := s.Service.VerifiyOTP(req.Email, req.OTP); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "email verified successfully",
	})
}

//CompleteSignup
func (s *AuthController) CompleteSignup(c *fiber.Ctx) error {
	var req dto.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	if err := utils.Validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := s.Service.CompleteSignup(&req); err != nil {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Profile created successfully",
	})
}

// Login
func (s *AuthController) Login(c *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": utils.ErrorMessage(utils.BADREQUEST, err),
			"err":   "biding error",
		})
	}
	if err := utils.Validator.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid input",
			"err":   err.Error(),
		})
	}

	user, access, refresh, err := s.Service.Login(&req)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access",
		Value:    access,
		HTTPOnly: true,
		SameSite: "Lax", // Strict in production
		Secure:   false, // true in production (HTTPS)
		Path:     "/",
		MaxAge:   60 * 15, // 15 minutes
	})
	c.Cookie(&fiber.Cookie{
		Name:     "refresh",
		Value:    refresh,
		HTTPOnly: true,
		SameSite: "Lax",
		Secure:   false,
		Path:     "/",  // only sent to refresh endpoint
		MaxAge:   60 * 60 * 24 * 7, // 7 days
	})

	return c.JSON(fiber.Map{
		"message": "login successful",
		"user": fiber.Map{
			"id":       user.ID,
			"email":    user.Email,
			// "username": user.Username,
			"access":   access,
		},
	})

}

// Refresh func
func (s *AuthController) Refresh(c *fiber.Ctx) error {
	access, refresh, err := s.Service.Refresh(c.Cookies("refresh"))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	c.Cookie(&fiber.Cookie{Name: "access",Value: access, HTTPOnly: true, SameSite: "Lax", Secure: false, Path: "/", MaxAge: 60 * 15,})
	c.Cookie(&fiber.Cookie{Name: "refresh", Value: refresh, HTTPOnly: true, SameSite: "Lax", Secure:false, Path: "/", MaxAge: 60 * 60 * 24 * 7,})

	return c.JSON(fiber.Map{
		"message": "Rotated",
	})
}

// Logout
func (s *AuthController) Logout(c *fiber.Ctx) error {
	s.Service.Logout(c.Cookies("access"), c.Cookies("refresh"))
	c.ClearCookie("access")
	c.ClearCookie("refresh")

	return c.JSON(fiber.Map{
		"message": "logout successfull",
	})
}

//Profile
// func (s *AuthController) Profile(c *fiber.Ctx) error {
// 	email := c.Locals("email").(string)

// 	user, err :=s.Service.Profile(email)
// 	if err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
// 			"error": err.Error(),
// 			"err" : "profile not showing",
// 		})
// 	}

// 	return c.Status(200).JSON(fiber.Map{
// 		"username" : user.Username,
// 		"email" : user.Email,
// 	})
// }
