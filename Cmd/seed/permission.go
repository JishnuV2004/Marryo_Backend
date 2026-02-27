package seed

import (
	models "marryo/Internal/Models"

	"gorm.io/gorm"
)

func SeedPermissions(db *gorm.DB) {

	permissions := []models.Permission{
		{Name: "ADMIN_DASHBOARD_VIEW"},
		{Name: "ADMIN_USER_MANAGEMENT_VIEW"},
		{Name: "ADMIN_MATCHES_VIEW"},
		{Name: "ADMIN_TELECALLING_VIEW"},
		{Name: "ADMIN_WEB_CONFIGURATION_VIEW"},
		{Name: "USER_HOME_VIEW"},
		{Name: "USER_MATCHES_VIEW"},
		{Name: "USER_INTERESTS_VIEW"},
		{Name: "USER_MESSAGES_VIEW"},
		{Name: "USER_SEARCH_VIEW"},
		{Name: "USER_NOTIFICATION_VIEW"},
	}

	for _, permission := range permissions {
		db.FirstOrCreate(&permission, models.Permission{Name: permission.Name})
	}
}

func SeedAdminRole(db *gorm.DB) error {

	var adminRole models.Role
	if err := db.Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var permissions []models.Permission
	if err := db.Find(&permissions).Error; err != nil {
		return err
	}

	// Replace all permissions for admin
	if err := db.Model(&adminRole).
		Association("Permissions").
		Replace(&permissions); err != nil {
		return err
	}

	return nil
}