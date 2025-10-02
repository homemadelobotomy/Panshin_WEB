package handler

import (
	"lab/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandlers(router *gin.Engine) {
	router.GET("/panels", h.GetSolarPanels)
	router.GET("/panels/:id", h.GetSolarPanel)
	router.GET("/solarpanel-requests/:id", h.GetSolarPanelRequest)
	router.POST("/solarpanel-requests/:solarpanel_id", h.AddSolarPanelToRequest)
	router.POST("/solarpanel-requests/delete/:request_id", h.DeleteSolarPanelRequest)

}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}
func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("templates/*")
	router.Static("/styles", "./resources/styles")
}
