// initialize package
package initialize

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"echo-learning/dao"
	"echo-learning/router"
	"echo-learning/services/mysql"
)

// init server.
func ServerInit() {
	// init mysql dependence.
	mysql.Init()

	// auto create table.
	AutoMigrate()

	// init echo instance.
	echoIns := GetDefaultFramework()

	// init routers.
	router.DispatchRouters(echoIns)

	// start server.
	ServerStart(echoIns)
}

// create echo instance.
func GetDefaultFramework() *echo.Echo {
	e := echo.New()
	e.HideBanner = true

	// Orderly, please do not disturb.
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.CORS())

	return e
}

// auto create table when not exist.
func AutoMigrate() {
	if err := mysql.GetDB().AutoMigrate(&dao.Course{}); err != nil {
		log.Fatalf("auto migrate failure, err=%v", err)
	}
}

// server start.
func ServerStart(echoIns *echo.Echo) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("server start and listen port=%s", port)
	if err := echoIns.Start(":" + port); err != nil {
		log.Fatalf("server start failure, err=%v", err)
	}
}
