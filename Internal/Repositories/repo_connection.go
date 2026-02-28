package repositories

import (
	models "marryo/Internal/Models"

	"gorm.io/gorm"
)

type PgSQLRepository struct {
	DB *gorm.DB
}

func NewRepo(db *gorm.DB) Repository{
	return &PgSQLRepository{DB: db}
}

//User Creation
func (r *PgSQLRepository) Create(req interface{}) error{
	return r.DB.Create(req).Error
}

func (r *PgSQLRepository) FindOne(model interface{}, query string, args ...any) error{
	return r.DB.Where(query, args...).First(model).Error
}

//Save
func (r *PgSQLRepository) Save(model interface{}) error {
	return r.DB.Save(model).Error
}

//GetUserByID

// func (r *PgSQLRepository) FindByID(model interface{}, Id uint, preloads ... string) error
// return db.First(model, Id).Error
func (r *PgSQLRepository) FindByID(model interface{}, Id uint, query string, preloads ... string) error {
	db :=  r.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db.Where(query, Id).First(model).Error
	// return db.First(model, Id).Error
}

//DeleteByID
func (r *PgSQLRepository) DeleteByID(model interface{}, id uint) error {
	return r.DB.Unscoped().Delete(model, id).Error
}

//update
func (r *PgSQLRepository) Update(model interface{}, fields map[string]interface{}, query string , args ...interface{}) error {
	return r.DB.Model(model).Where(query, args ...).Updates(fields).Error
}

//Count 
func (r *PgSQLRepository) CountWithCondition(model interface{}, query string, args ...interface{}) (int64, error) {
	var count int64
	err := r.DB.Model(model).Where(query, args...).Count(&count).Error
	return count, err
}


/////////// RBAC. ///////////

func (r *PgSQLRepository) UpdatePermissions(roleID uint, permissionIDs []uint) error {

	var role models.Role
	if err := r.DB.First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := r.DB.Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
		return err
	}

	return r.DB.Model(&role).Association("Permissions").Replace(&permissions)
}


func (r *PgSQLRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.DB.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *PgSQLRepository) GetUserWithRoles(userID uint) (*models.User, error) {

	var user models.User

	err := r.DB.
		Preload("Roles").
		Preload("Roles.Permissions").
		First(&user, userID).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}


func (r *PgSQLRepository) GetPermissionsByRole(role string) (*models.Role, error) {

	var roles models.Role

	err := r.DB.
		Preload("Permissions").
		Where("name = ?", role).
		First(&roles).Error

	if err != nil {
		return nil, err
	}
	

	return &roles, nil
}