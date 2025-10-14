package handler

import (
	"errors"
	dto "lab/internal/app/DTO"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (h *Handler) RegisterUserHandlers(router *gin.Engine) {
	router.POST("/api/user/registration", h.Registration)
	router.GET("/api/user/:id", h.GetUserData)
	router.PUT("/api/user", h.ChangeUserData)
}

func (h *Handler) Registration(ctx *gin.Context) {
	var user dto.UserRegistration
	var pgErr *pgconn.PgError
	if err := ctx.BindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}

	response, err := h.Service.AddNewUser(user)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"")
		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			h.errorHandler(ctx, http.StatusBadRequest,
				"пользователь с таким логином уже существует")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (h *Handler) GetUserData(ctx *gin.Context) {
	userId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}

	response, err := h.Service.GetUserData(uint(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)

}

func (h *Handler) ChangeUserData(ctx *gin.Context) {
	var user dto.ChangeUserData
	if err := ctx.BindJSON(&user); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}
	response, err := h.Service.ChangeUserData(user)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusOK, response)
}
