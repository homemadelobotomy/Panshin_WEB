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
	r.GET("/panels", handler.GetSolarpanels)
	r.GET("/panel/:id", handler.GetSolarpanel)
	r.GET("/bid/:id", handler.GetBid)

	r.Run()

	log.Println("Server down")
}
