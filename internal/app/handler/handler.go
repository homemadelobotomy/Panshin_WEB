package handler

import (
	"lab/internal/app/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
}

var STATUS_CODES = map[int]string{
	403: "Forbidden",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

func (h *Handler) RegisterHandlers(router *gin.Engine) {
	h.RegisterSolarPanelHandlers(router)
	h.RegisterSolarPanelsRequestHandlers(router)

}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, message string) {
	ctx.JSON(errorStatusCode, gin.H{
		"status":  strconv.Itoa(errorStatusCode) + " " + STATUS_CODES[errorStatusCode],
		"message": message,
	})
}
