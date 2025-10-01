package main

import (
	"fmt"

	"lab/internal/app/config"
	"lab/internal/app/dsn"
	"lab/internal/app/handler"
	"lab/internal/app/repository"
	"lab/internal/app/service"
	"lab/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}
	serv := service.NewService(rep)

	hand := handler.NewHandler(serv)

	SolarPanelPowerCalculator := pkg.NewApp(conf, router, hand)
	SolarPanelPowerCalculator.RunApp()
}
