package models

import (
	"gorm.io/gorm"
)

type Scope struct {
	gorm.Model
	ScopeName string `gorm:"unique" json:"scope_name"`
}
