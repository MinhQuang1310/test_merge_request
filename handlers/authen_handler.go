package handlers

import (
	"fmt"
	"net/http"
	"test-echo/auth"
	"test-echo/common"
	"test-echo/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Login(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Lấy thông tin đăng nhập từ request và gán vào input
		var input models.User
		if err := c.Bind(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		//Validate dữ liệu
		if input.UserName == "" || input.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Username and Password must not be empty"})
		}

		//Kiểm tra username có tồn tại không
		var user models.User
		if result := db.Where("user_name = ?", input.UserName).First(&user); result.Error != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
		}

		//Kiểm tra password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, map[string]string{"message": "Invalid username or password"})
		}

		//Lấy ra các quyền của user đó
		var userRoles []models.UserRole
		if err := db.Where("user_id = ?", user.ID).Find(&userRoles).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Error retrieving user roles"})
		}

		// Tạo mảng chứa tên các role
		var roleNames []string
		for _, ur := range userRoles {
			var role models.Role
			if err := db.First(&role, ur.RoleID).Error; err == nil {
				roleNames = append(roleNames, role.RoleName)
			}
		}

		// Tạo token với các role
		token, err := auth.CreateToken(user.ID, roleNames...)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "An error occurred while signing the token"})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Login successful!", "token": token})
	}
}

func Register(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.User
		if err := c.Bind(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		if input.Password == "" || input.UserName == "" || input.Email == "" {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Username, Email and Password must not be empty"})
		}

		//Mã hóa password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Could not hash the password"})
		}

		user := models.User{
			UserName: input.UserName,
			Email:    input.Email,
			Password: string(hashedPassword),
		}

		result := db.Where(models.User{Email: input.Email}).FirstOrCreate(&user)

		if result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": result.Error.Error()})
		}

		if result.RowsAffected > 0 {
			userRoleUser := models.UserRole{UserID: user.ID, RoleID: 2}
			userRoleViewer := models.UserRole{UserID: user.ID, RoleID: 3}
			db.Create(&userRoleUser)
			db.Create(&userRoleViewer)

			return c.JSON(http.StatusOK, map[string]interface{}{
				"message": "User created",
				"user":    user,
			})
		} else {
			return echo.NewHTTPError(http.StatusConflict, map[string]string{"message": "Email already exists"})
		}
	}
}

func RequestResetPassword(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(struct {
			Email string `json:"email"`
		})
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		// Kiểm tra xem email có tồn tại trong cơ sở dữ liệu không
		var user models.User
		if err := db.Where("email = ?", req.Email).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Email does not exist"})
		}

		// Tạo token đặt lại mật khẩu
		token, err := auth.CreateResetPassowrdToken(user.ID)

		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		// Gửi email với token
		common.SendEmail(req.Email, token)

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "A reset password email has been sent",
		})
	}
}

func Hello() {
	fmt.Println("Hello")
}

func ResetPassword(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.QueryParam("token")

		// Kiểm tra token
		claims, err := auth.ValidateResetPasswordToken(tokenString)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid or expired token"})
		}

		req := new(struct {
			Password string `json:"password"`
		})
		if err := c.Bind(req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Invalid input"})
		}

		var user models.User
		if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "Email does not exist"})
		}

		// Mã hóa mật khẩu mới
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": "Could not hash the password"})
		}

		user.Password = string(hashedPassword)

		if result := db.Save(&user); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, result.Error.Error())
		}

		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Password has been reset",
		})
	}
}

func AuthenGroup(e *echo.Group, db *gorm.DB) *echo.Group {
	e.GET("/login", Login(db))
	e.POST("/register", Register(db))
	e.GET("/reset-password", RequestResetPassword(db))
	e.PUT("/reset-password", ResetPassword(db))
	return e
}
