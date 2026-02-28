package controller

import (
	"log"
	dto "marryo/Internal/DTO"
	services "marryo/Internal/Services"
	"strconv"

	"github.com/gofiber/fiber/v2"
)


type RoleHandler struct {
	service services.RoleService
}

func NewRoleHandler(service services.RoleService) *RoleHandler {
	return &RoleHandler{service}
}


//create role
func (h *RoleHandler) CreateRole(c *fiber.Ctx) error {

	var req dto.CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.CreateRole(req.Name); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Role created successfully"})
}

//update role
func (h *RoleHandler) UpdateRolePermissions(c *fiber.Ctx) error {

	var req dto.UpdateRolePermissionsRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := h.service.UpdateRolePermissions(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Permissions updated"})
}

//get role permissions
func (h *RoleHandler) GetRolePermissions(c *fiber.Ctx) error {

	roleID, _ := strconv.Atoi(c.Params("id"))
    log.Println("id",roleID)
	permissions, err := h.service.GetRolePermissions(uint(roleID))
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(permissions)
}

func (h *RoleHandler) GetAllPermissions(c *fiber.Ctx) error {

	permissions, err := h.service.GetAllRoles()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(permissions)
}