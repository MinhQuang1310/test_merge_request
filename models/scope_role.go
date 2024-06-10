package models

type ScopeRole struct {
	ScopeID uint `gorm:"primaryKey"`
	RoleID  uint `gorm:"primaryKey"`
}
