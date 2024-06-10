package migrations

import (
	"gorm.io/gorm"
)

func Migration(gormDB *gorm.DB) {
	MigratePermission(gormDB)
	MigrateRole(gormDB)
	MigrateScope(gormDB)
	MigrateUser(gormDB)
	MigrateScope(gormDB)
	MigrateUserRole(gormDB)
	MigrateRolePermission(gormDB)
	MigrateScopeRole(gormDB)
	MigrateBlog(gormDB)
}
