package handler

import (
	"lab/internal/app/repository"

	"metoda/internal/app/repository"

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

func (h *Handler) RegistrHandlers(router *gin.Engine) {
	router.GET("/", h.GetSolarPanels)
	router.GET("/panel/:id", h.GetSolarPanel)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

func (h *Handler) GetSolarPanels(ctx *gin.Context) {

}

func (h *Handler) GetSolarPanel(ctx *gin.Context) {

}
