package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateScopeRole(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.ScopeRole{}) {
		gormDB.AutoMigrate(&models.ScopeRole{})
		//Scope for ADMIN
		gormDB.Create(&models.ScopeRole{ScopeID: 2, RoleID: 1}) //ALL USERS
		gormDB.Create(&models.ScopeRole{ScopeID: 4, RoleID: 1}) //ALL BLOGS
		//Scope for User
		gormDB.Create(&models.ScopeRole{ScopeID: 2, RoleID: 2}) //ALL USERS
		gormDB.Create(&models.ScopeRole{ScopeID: 4, RoleID: 2}) //ALL BLOGS
		//Scope for AUTHOR
		gormDB.Create(&models.ScopeRole{ScopeID: 1, RoleID: 3}) //OWNED USER
		gormDB.Create(&models.ScopeRole{ScopeID: 3, RoleID: 3}) //OWNED BLOGS
	}
}
