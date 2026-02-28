package repositories

import models "marryo/Internal/Models"

type Repository interface {
	Create(req interface{}) error
	FindOne(model interface{}, query string, args ...any) error
	Save(model interface{}) error
	FindByID(model interface{}, Id uint, query string, preloads ... string) error
	Update(model interface{}, fields map[string]interface{}, query string , args ...interface{}) error
	CountWithCondition(model interface{}, query string, args ...interface{}) (int64, error)
	// FindByIDWithPreload(model interface{}, Id uint) error
	DeleteByID(model interface{}, id uint) error

	UpdatePermissions(roleID uint, permissionIDs []uint) error 
	GetAll() ([]models.Role, error)
	GetUserWithRoles(userID uint) (*models.User, error)

	GetPermissionsByRole(role string) (*models.Role, error)
}