package main

import (
	"test-echo/common"
	"test-echo/db"
	"test-echo/handlers"
	"test-echo/migrations"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {

	//common.SendEmail()

	gormDB := db.InitDB()
	sqlDB, err := gormDB.DB()

	migrations.Migration(gormDB)

	if err != nil {
		panic("failed to get sql.DB")
	}
	defer sqlDB.Close()

	e := echo.New()

	// e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
	// 	Format: "Time=${time_rfc3339}, Method=${method}, Uri=${uri}, Status=${status}\n",
	// }))
	// e.Use(middleware.Logger())
	common.NewLogger()
	e.Use(common.LoggingMiddleware)
	e.Use(middleware.Recover())

	//Route API
	api := e.Group("/api")
	v1 := api.Group("/v1")
	handlers.UsersGroup(v1, gormDB)
	handlers.BlogsGroup(v1, gormDB)
	handlers.AuthenGroup(v1, gormDB)

	e.Start(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
