package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateRole(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.Role{}) {
		gormDB.AutoMigrate(&models.Role{})
		gormDB.Create(&models.Role{RoleName: "ADMIN"})
		gormDB.Create(&models.Role{RoleName: "USER"})
		gormDB.Create(&models.Role{RoleName: "AUTHOR"})
	}
}
