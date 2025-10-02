package handler

import (
	"errors"
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"lab/internal/app/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// import "github.com/gin-gonic/gin"

// func (h *Handler) RegisterSolarRequestHandlers(router *gin.Engine) {
// 	router.GET("/request/:id", h.GetSolarPanelRequest)
// 	router.POST("/request.add/:solarpanel_id", h.AddSolarPanelToRequest)
// 	router.POST("/request/delete/:request_id", h.DeleteSolarPanelRequest)
// }

func (h *Handler) RegisterRequestPanelsHandlers(router *gin.Engine) {
	router.DELETE("/api/solarpanel-requests/:id/:solarpanelId", h.DeleteSolarPanelFromRequest)
	router.PUT("/api/solarpanel-requests/:id/:solarpanelId", h.ChangeSolarPanelArea)
}

func (h *Handler) DeleteSolarPanelFromRequest(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	solarPanelId, err := strconv.Atoi(ctx.Param("solarpanelId"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}

	err = h.Service.DeleteSolarPanelFromRequest(ds.GetUser().GetId(), uint(requestId), uint(solarPanelId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else if errors.Is(err, service.ErrNoRecords) {
			h.errorHandler(ctx, http.StatusForbidden,
				"солнечная панель с таким id отсутсвует в заявке")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200 OK",
		"message": "солнечная панель успешно удалена из заявки по солнечным панелям",
	})
}

func (h *Handler) ChangeSolarPanelArea(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	solarPanelId, err := strconv.Atoi(ctx.Param("solarpanelId"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}

	var area dto.ChangeSolarPanelAreaRequest

	if err = ctx.BindJSON(&area); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}
	response, err := h.Service.ChangeSolarPanelAreaInRequest(ds.GetUser().GetId(), uint(requestId), uint(solarPanelId), area.Area)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else if errors.Is(err, service.ErrNoRecords) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечная панель с таким id отсутсвует в заявке")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено неверное значение для площади")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}
