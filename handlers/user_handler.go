package handlers

import (
	"net/http"
	"strconv"
	"test-echo/auth"
	"test-echo/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func GetUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if result := db.First(&user, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		return c.JSON(http.StatusOK, user)
	}
}

func GetUsers(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var users []models.User
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit == 0 {
			limit = 20
		}
		offset := (page - 1) * limit

		if result := db.Limit(limit).Offset(offset).Find(&users); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		return c.JSON(http.StatusOK, users)
	}
}

func UpdateUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var user models.User
		if result := db.First(&user, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}

		if err := c.Bind(&user); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		// Kiểm tra xem mật khẩu đã thay đổi hay không. Nếu có thì mã hóa và thay đổi
		if user.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError, "Could not hash the password")
			}
			user.Password = string(hashedPassword)
		}

		if result := db.Save(&user); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "User updated",
			"user":    user,
		})
	}
}

func DeleteUser(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if result := db.Delete(&models.User{}, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "User deleted"})
	}
}

func UsersGroup(e *echo.Group, db *gorm.DB) *echo.Group {
	usersGroup := e.Group("/users")

	usersGroup.Use(auth.CheckRole(db, "ADMIN", "USER", "AUTHOR"))

	usersGroup.GET("/:id", GetUser(db))
	usersGroup.GET("/", GetUsers(db))
	usersGroup.PUT("/:id", UpdateUser(db))
	usersGroup.DELETE("/:id", DeleteUser(db))

	return usersGroup
}
