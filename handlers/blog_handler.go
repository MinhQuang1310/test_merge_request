package handlers

import (
	"net/http"
	"strconv"
	"test-echo/auth"
	"test-echo/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CreateBlog(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var input models.Blog
		if err := c.Bind(&input); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"message": "invalid input"})
		}
		claims, err := auth.GetClaimsFromToken(c)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"message": err.Error()})
		}

		blog := models.Blog{
			Title:   input.Title,
			Content: input.Content,
			UserID:  claims.UserID,
		}

		if result := db.Create(&blog); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "could not create blog"})
		}

		return c.JSON(http.StatusOK, blog)
	}
}

func GetBlog(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var blog models.Blog
		if result := db.First(&blog, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Blog not found")
		}

		return c.JSON(http.StatusOK, blog)
	}
}

func GetBlogs(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var blogs []models.Blog
		page, _ := strconv.Atoi(c.QueryParam("page"))
		if page == 0 {
			page = 1
		}
		limit, _ := strconv.Atoi(c.QueryParam("limit"))
		if limit == 0 {
			limit = 20
		}
		offset := (page - 1) * limit

		if result := db.Limit(limit).Offset(offset).Find(&blogs); result.Error != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return c.JSON(http.StatusOK, blogs)
	}
}

func UpdateBlog(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var blog models.Blog

		if result := db.First(&blog, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusNotFound, map[string]string{"error": "Blog not found"})
		}

		if err := c.Bind(&blog); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		if result := db.Save(&blog); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return c.JSON(http.StatusOK, blog)
	}
}

func DeleteBlog(db *gorm.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		if result := db.Delete(&models.Blog{}, c.Param("id")); result.Error != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, map[string]string{"error": result.Error.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{"message": "Blog deleted"})
	}
}

func BlogsGroup(e *echo.Group, db *gorm.DB) *echo.Group {
	blogsGroup := e.Group("/blogs")

	blogsGroup.Use(auth.CheckRole(db, "ADMIN", "USER"))

	blogsGroup.POST("/", CreateBlog(db))
	blogsGroup.GET("/:id", GetBlog(db))
	blogsGroup.GET("/", GetBlogs(db))
	blogsGroup.PUT("/:id", UpdateBlog(db))
	blogsGroup.DELETE("/:id", DeleteBlog(db))

	return blogsGroup
}
