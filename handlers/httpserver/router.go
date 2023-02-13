package httpserver

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/vano2903/service-template/controller"

	_ "github.com/vano2903/service-template/docs"
)

//The package is called httpserver and not http because it's a bad practice
//to name a package after a standard library package (net/http in this case)

// Swagger spec:
//	@title			Go Service Template
//	@version		1.0
//	@description	User Management Service
//	@contact.name	Vano2903
//	@contact.url	https://github.com/vano2903
//	@contact.email	davidevanoncini2003@gmail.com
//	@host			localhost:8080
//	@BasePath		/api/v1
func InitRouter(e *echo.Echo, l *logrus.Logger, userController *controller.User) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	api := e.Group("/api/v1")
	user := api.Group("/user")

	//user routes

	userHttpHandler := NewUserHttpHandler(user, userController, l)
	userHttpHandler.RegisterRoutes()
}
