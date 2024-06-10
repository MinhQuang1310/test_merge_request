package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateBlog(gormDB *gorm.DB) {
	// Kiểm tra xem bảng Permission đã tồn tại hay chưa
	if !gormDB.Migrator().HasTable(&models.Blog{}) {
		// Nếu chưa tồn tại, thực hiện migration
		gormDB.AutoMigrate(&models.Blog{})
	}
}
