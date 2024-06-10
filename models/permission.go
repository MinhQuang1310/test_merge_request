package models

import (
	"gorm.io/gorm"
)

type Permission struct {
	gorm.Model
	PermissionName string `gorm:"unique" json:"permission_name"`
}
