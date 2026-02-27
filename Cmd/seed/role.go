package seed

import (
	models "marryo/Internal/Models"

	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) error {

	roles := []models.Role{
		{Name: "admin"},
		{Name: "staff"},
		{Name: "user"},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, models.Role{Name: role.Name}).Error; err != nil {
			return err
		}
	}

	return nil
}