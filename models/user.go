package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName string `json:"username"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Password string `json:"password"`
}
