package repositories

import (

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
func (r *PgSQLRepository) FindByID(model interface{}, Id uint, preloads ... string) error {
	db :=  r.DB
	for _, p := range preloads {
		db = db.Preload(p)
	}
	return db.First(model, Id).Error
}

// func (r *PgSQLRepository) FindByIDWithPreload(model interface{}, Id uint) error {
// 	return r.DB.Preload("profile").First(model, Id).Error
// }

//DeleteByID
func (r *PgSQLRepository) DeleteByID(model interface{}, id uint) error {
	return r.DB.Unscoped().Delete(model, id).Error
}