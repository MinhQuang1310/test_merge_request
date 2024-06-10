package auth

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"test-echo/common"
	"test-echo/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CheckRole(gormDB *gorm.DB, allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Lấy token từ header Authorization
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return echo.NewHTTPError(http.StatusUnauthorized, "Missing Authorization header")
			}

			tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

			// Xác thực token
			claims, checkErr := validateToken(tokenString)
			if checkErr != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, checkErr.Error())
			}

			// Kiểm tra quyền truy cập
			if len(allowedRoles) > 0 {
				hasAccess := false
				for _, role := range allowedRoles {
					if claims.hasRole(role) {
						hasAccess = true
						break
					}
				}

				if !hasAccess {
					return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Role Unauthorized"})
				}

				//Lấy role ID
				currentRoleIDs, checkErr := getRoleIDs(gormDB, claims.Roles...)
				if checkErr != nil {
					return echo.NewHTTPError(http.StatusForbidden, checkErr)
				}

				//Kiểm tra từng role
				for _, currentRoleID := range currentRoleIDs {
					//Kiểm tra scope
					if err := checkScope(gormDB, c, currentRoleID, claims.UserID); err != nil {
						return err
					}
					//Kiểm tra permission
					if err := checkPermission(gormDB, c, currentRoleID); err != nil {
						return err
					}
				}
			}

			// Lưu claims vào context
			c.Set("user", claims)

			// Gọi hàm tiếp theo trong chuỗi middleware
			return next(c)
		}
	}
}

func checkScope(gormDB *gorm.DB, c echo.Context, roleID uint, userID uint) error {
	currentScopes, checkErr := getScopeNames(gormDB, roleID)
	if checkErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, checkErr)
	}
	switch currentGroup(c) {
	case "users":
		if contains(currentScopes, "ALL_USERS") || contains(currentScopes, "OWNED_USER") {
			if err := checkOwned(userID, gormDB, c, currentScopes...); err != nil {
				return err
			}
		}
	case "blogs":
		if contains(currentScopes, "ALL_BLOGS") || contains(currentScopes, "OWNED_BLOGS") {
			if err := checkOwned(userID, gormDB, c, currentScopes...); err != nil {
				return err
			}
		}
	default:
		return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Can not get group"})
	}
	return nil
}

func checkOwned(userID uint, gormDB *gorm.DB, c echo.Context, scopes ...string) error {
	if c.Request().Method == echo.GET {
		return nil
	}
	for _, scope := range scopes {
		switch scope {
		case "OWNED_USER":
			id, err := strconv.Atoi(c.Param("id")) //Lấy id từ param và đổi về int
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
			}
			if userID != uint(id) {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Scope Unauthorized"})
			}
			return nil
		case "OWNED_BLOGS":
			id, err := strconv.Atoi(c.Param("id")) //Lấy id từ param và đổi về int
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid ID"})
			}
			var blog models.Blog
			if err := gormDB.Where("id = ?", id).Find(&blog).Error; err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Cannot find Blogs"})
			}
			if userID != blog.UserID {
				return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Scope Unauthorized"})
			}
			return nil
		default:
			return nil
		}
	}
	return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Not have any scope"})
}

func checkPermission(gormDB *gorm.DB, c echo.Context, roleID uint) error {
	currentPermissions, checkErr := getPermissionNames(gormDB, roleID)
	if checkErr != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, checkErr)
	}

	switch c.Request().Method {
	case echo.POST:
		if contains(currentPermissions, "CREATE") {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "You don't have permission"})
	case echo.GET:
		if contains(currentPermissions, "READ") {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "You don't have permission"})
	case echo.PUT:
		if contains(currentPermissions, "UPDATE") {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "You don't have permission"})
	case echo.DELETE:
		if contains(currentPermissions, "DELETE") {
			return nil
		}
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "You don't have permission"})
	default:
		return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "You don't have permission"})
	}
}

func getPermissionNames(gormDB *gorm.DB, roleID uint) ([]string, error) {
	var permissions []models.Permission

	if err := gormDB.Joins("JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Where("role_permissions.role_id = ?", roleID).
		Find(&permissions).Error; err != nil {
		return nil, err
	}

	var permissionNames []string
	for _, permission := range permissions {
		permissionNames = append(permissionNames, permission.PermissionName)
	}

	return permissionNames, nil
}

func currentGroup(c echo.Context) string {
	path := c.Path()
	parts := strings.Split(path, "/")
	if len(parts) > 3 {
		return parts[3] // Trả về 'blogs' hoặc 'users' tùy thuộc vào đường dẫn
	}
	return ""
}

func getScopeNames(gormDB *gorm.DB, roleID uint) ([]string, error) {
	var scopes []models.Scope

	if err := gormDB.Joins("JOIN scope_roles ON scopes.id = scope_roles.scope_id").
		Where("scope_roles.role_id = ?", roleID).
		Find(&scopes).Error; err != nil {
		return nil, err
	}

	var scopeNames []string
	for _, scope := range scopes {
		scopeNames = append(scopeNames, scope.ScopeName)
	}

	return scopeNames, nil
}

func getRoleIDs(gormDB *gorm.DB, roleNames ...string) ([]uint, error) {
	var roles []models.Role
	if err := gormDB.Where("role_name IN ?", roleNames).Find(&roles).Error; err != nil {
		return nil, err
	}

	var roleIDs []uint
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
	}

	return roleIDs, nil
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func (c *JwtCustomClaims) hasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}

func validateToken(tokenString string) (*JwtCustomClaims, error) {
	common.LoadEnvFile()
	sKey := os.Getenv("JWT_KEY")

	token, err := jwt.ParseWithClaims(tokenString, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method=%v", token.Header["alg"])
		}
		return []byte(sKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
