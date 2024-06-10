package migrations

import (
	"test-echo/models"

	"gorm.io/gorm"
)

func MigrateUser(gormDB *gorm.DB) {
	if !gormDB.Migrator().HasTable(&models.User{}) {
		gormDB.AutoMigrate(&models.User{})
		gormDB.Create(&models.User{UserName: "Lam Chi Tinh", Email: "lamchitinh@gmail.com", Password: "$2a$10$VwWrgovBQvetRpW9eRS6zuJDQ4ALAXPiqn/YqSn8X8s75pXQ73fEy"}) //Pass: Lam Chi
		gormDB.Create(&models.User{UserName: "User 1", Email: "User1@gmail.com", Password: "$2a$10$WsXlCmzPH9xSvdSg9J2xZONs6zmOtrMVTqzgoG8ekeltIVCPaNJpm"})            //Pass: User1
	}
}
