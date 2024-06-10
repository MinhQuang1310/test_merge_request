package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateScope(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.Scope{}) {
		gormDB.AutoMigrate(&models.Scope{})
		gormDB.Create(&models.Scope{ScopeName: "OWNED_USER"})
		gormDB.Create(&models.Scope{ScopeName: "ALL_USERS"})
		gormDB.Create(&models.Scope{ScopeName: "OWNED_BLOGS"})
		gormDB.Create(&models.Scope{ScopeName: "ALL_BLOGS"})
	}
}
