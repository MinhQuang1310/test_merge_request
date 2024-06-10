package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateUserRole(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.UserRole{}) {
		gormDB.AutoMigrate(&models.UserRole{})
		gormDB.Create(&models.UserRole{UserID: 1, RoleID: 1})
		gormDB.Create(&models.UserRole{UserID: 2, RoleID: 2})
		gormDB.Create(&models.UserRole{UserID: 2, RoleID: 3})
	}
}
