package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigratePermission(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.Permission{}) {
		gormDB.AutoMigrate(&models.Permission{})
		gormDB.Create(&models.Permission{PermissionName: "CREATE"})
		gormDB.Create(&models.Permission{PermissionName: "READ"})
		gormDB.Create(&models.Permission{PermissionName: "UPDATE"})
		gormDB.Create(&models.Permission{PermissionName: "DELETE"})
	}
}
