package api

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"Project/internal/app/handler"
	"Project/internal/app/repository"
)

func StartServer() {
	log.Println("Server start up")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("repository init error")
	}
	handler := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/resources", "./resources")
	r.GET("/panels", handler.GetSolarPanels)
	r.GET("/panel/:id", handler.GetSolarPanel)
	r.GET("/solar_panels_request/:id", handler.GetBid)

	r.Run("127.0.0.1:8001")

	log.Println("Server down")
}
