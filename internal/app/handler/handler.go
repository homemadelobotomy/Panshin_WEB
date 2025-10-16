package handler

import (
	"lab/internal/app/config"
	"lab/internal/app/service"
	"lab/internal/redis"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
	Redis   *redis.Client
	Config  *config.Config
}

var STATUS_CODES = map[int]string{
	403: "Forbidden",
	400: "Bad Request",
	404: "Not Found",
	500: "Internal Server Error",
}

func NewHandler(s *service.Service, redis *redis.Client, c *config.Config) *Handler {
	return &Handler{
		Service: s,
		Redis:   redis,
		Config:  c,
	}
}

func (h *Handler) RegisterHandlers(router *gin.Engine) {
	h.RegisterSolarPanelHandlers(router)
	h.RegisterSolarPanelsRequestHandlers(router)
	h.RegisterRequestPanelsHandlers(router)
	h.RegisterUserHandlers(router)
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, message string) {
	ctx.JSON(errorStatusCode, gin.H{
		"status":  strconv.Itoa(errorStatusCode) + " " + STATUS_CODES[errorStatusCode],
		"message": message,
	})
}
