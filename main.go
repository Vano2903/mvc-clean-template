package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/vano2903/service-template/config"
	"github.com/vano2903/service-template/controller"
	"github.com/vano2903/service-template/handlers/httpserver"
	"github.com/vano2903/service-template/model"
	"github.com/vano2903/service-template/pkg/logger"
	"github.com/vano2903/service-template/providers/logo"
	"github.com/vano2903/service-template/repo/mock"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	l := logger.NewLogger(conf.Log.Level, conf.Log.Type)
	l.Debug("initizalized logger")

	if conf.Database.Driver != "mock" {
		log.Fatal("only mock database is supported in this example")
	}

	//creating the instances for the application
	repo := mock.NewRepo()
	logoService := logo.NewServiceLogo(conf.Services.Logo.ApiKey, conf.Services.Logo.BaseUrl)

	//creating the controller
	c := controller.NewUserController(repo, logoService, l)

	GenerateExampleEntries(l, c)

	//creating the http server
	e := echo.New()
	httpserver.InitRouter(e, l, c, conf)

	//starting the server
	e.Logger.Fatal(e.Start(":" + conf.HTTP.Port))
}

func GenerateExampleEntries(l *logrus.Logger, c *controller.User) {
	if _, err := c.CreateUser("Davide", "Vanoncini", "davidevanoncini2003@gmail.com", "password", model.RoleAdmin); err != nil {
		l.Fatal("unable to create test user")
	}
	if _, err := c.CreateUser("John", "Doe", "johndoe@bingchilling.cn", "123secure", model.RoleUser); err != nil {
		l.Fatal("unable to create test user")
	}
	if _, err := c.CreateUser("Foo", "Bar", "foo@bar.com", "psw1", model.RoleUnupdatable); err != nil {
		l.Fatal("unable to create test user")
	}
}
