package services

import (
	dto "marryo/Internal/DTO"
	models "marryo/Internal/Models"
	repositories "marryo/Internal/Repositories"
)


type RoleService struct {
	repo repositories.Repository
}

func NewRoleService(repo repositories.Repository) *RoleService {
	return &RoleService{repo : repo}
}


func (s *RoleService) CreateRole(name string) error {

	role := models.Role{
		Name: name,
	}

	return s.repo.Create(&role)
}

func (s *RoleService) GetRolePermissions(roleID uint) ([]string, error) {


	var role models.Role
	err := s.repo.FindByID(&role, roleID,  "id = ?", "Permissions")
	if err != nil {
		return nil, err
	}

	var permissions []string
	for _, p := range role.Permissions {
		permissions = append(permissions, p.Name)
	}

	return permissions, nil
}

// for update
func (s *RoleService) UpdateRolePermissions(req dto.UpdateRolePermissionsRequest) error {
	return s.repo.UpdatePermissions(req.RoleID, req.PermissionIDs)
}

func (s *RoleService) GetAllRoles() ([]models.Role, error) {
	return s.repo.GetAll()
}