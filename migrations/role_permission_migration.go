package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateRolePermission(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.RolePermission{}) {
		gormDB.AutoMigrate(&models.RolePermission{})
		//Permission for ADMIN
		gormDB.Create(&models.RolePermission{RoleID: 1, PermissionID: 1}) //CREATE
		gormDB.Create(&models.RolePermission{RoleID: 1, PermissionID: 2}) //READ
		gormDB.Create(&models.RolePermission{RoleID: 1, PermissionID: 3}) //UPDATE
		gormDB.Create(&models.RolePermission{RoleID: 1, PermissionID: 4}) //DELETE

		//Permission for USER
		gormDB.Create(&models.RolePermission{RoleID: 2, PermissionID: 1})
		gormDB.Create(&models.RolePermission{RoleID: 2, PermissionID: 2})

		// Permission for AUTHOR
		gormDB.Create(&models.RolePermission{RoleID: 3, PermissionID: 2})
		gormDB.Create(&models.RolePermission{RoleID: 3, PermissionID: 3})
		gormDB.Create(&models.RolePermission{RoleID: 3, PermissionID: 4})
	}
}
	