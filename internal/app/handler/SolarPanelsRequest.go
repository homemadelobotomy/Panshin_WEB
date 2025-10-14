package handler

import (
	"errors"
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"lab/internal/app/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h *Handler) RegisterSolarPanelsRequestHandlers(router *gin.Engine) {
	router.GET("/api/solarpanel-request-info", h.GetSolarPanelsInRequest)
	solarPanelRequestGroups := router.Group("/api/solarpanel-requests")
	{
		solarPanelRequestGroups.GET("", h.GetFilteredSolarPanelRequests)
		solarPanelRequestGroups.GET("/:id", h.GetOneSolarPanelRequest)
		solarPanelRequestGroups.PUT("/:id", h.ChangeSolarPanelRequest)
		solarPanelRequestGroups.PUT("/:id/formate", h.FormateSolarPanelRequest)
		solarPanelRequestGroups.PUT("/:id/moderate", h.ModeratorAction)
		solarPanelRequestGroups.DELETE("/:id", h.DeleteSolarPanelRequest)
	}
}

func (h *Handler) GetSolarPanelsInRequest(ctx *gin.Context) {
	userId := ds.GetUser().GetId()

	requestId, numberOfPanels, err := h.Service.GetSolarPanelsInRequest(userId)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusOK, dto.NumberOfPanelsResponse{
				RequestId:      0,
				NumberOfPanels: 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError,
				err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, dto.NumberOfPanelsResponse{
		RequestId:      requestId,
		NumberOfPanels: numberOfPanels,
	})
}

func (h *Handler) GetFilteredSolarPanelRequests(ctx *gin.Context) {
	var (
		startDate time.Time
		endDate   time.Time
		err       error
	)

	userId := ds.GetUser().GetId()

	start_date := ctx.Query("start_date")
	end_date := ctx.Query("end_date")
	status := ctx.Query("status")

	layout := "02-01-2006 15:04:05"

	if start_date != "" {
		startDate, err = time.Parse(layout, start_date)
		if err != nil {

			h.errorHandler(ctx, http.StatusBadRequest,
				"введен неверный формат start_date, правильный формат: dd-mm-yyyy hh:mm:ss")
			return
		}
	} else {
		startDate = time.Time{}
	}

	if end_date != "" {
		endDate, err = time.Parse(layout, end_date)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введен неверный формат end_date, правильный формат: dd-mm-yyyy hh:mm:ss")
			return
		}
	} else {
		endDate = time.Time{}
	}

	filter := dto.SolarPanleRequestFilter{
		Status:     status,
		Start_date: startDate,
		End_date:   endDate,
	}
	response, err := h.Service.GetFilteredSolarPanelRequests(userId, filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) || errors.Is(err, service.ErrNoRecords) {
			h.errorHandler(ctx, http.StatusNotFound,
				"для данного пользователя не найдено ни одной заявки по солнечным панелям")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) GetOneSolarPanelRequest(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
	}
	response, err := h.Service.GetOneSolarPanelRequest(uint(requestId), ds.GetUser().GetId())

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найден")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)

}

func (h *Handler) ChangeSolarPanelRequest(ctx *gin.Context) {
	var (
		insolationRequest dto.ChangeSolarPanelRequest
		response          dto.OneSolarPanelRequestResponse
	)
	if err := ctx.BindJSON(&insolationRequest); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	response, err = h.Service.ChangeSolarPanelRequest(ds.GetUser().GetId(), uint(requestId), insolationRequest.Insolation)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено некорректное значение для инсоляции")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) FormateSolarPanelRequest(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	response, err := h.Service.FormateSolarPanelRequest(uint(requestId), ds.GetUser().GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено некорректное значение для инсоляции или значения площади  < 0")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) ModeratorAction(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	var action dto.ModeratorAction
	if err = ctx.BindJSON(&action); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}
	response, err := h.Service.ModeratorAction(uint(requestId), action.Action, ds.GetUser().GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"пользователь не является модератором")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"действие модератора должно быть 'отклонен|завершен'")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSolarPanelRequest(ctx *gin.Context) {
	requestId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id заявки по солнечным панелям")
		return
	}
	err = h.Service.DeleteSolarPanelRequest(uint(requestId), ds.GetUser().GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найдено")
		} else if errors.Is(err, service.ErrForbidden) {
			h.errorHandler(ctx, http.StatusForbidden,
				"заявка по солнечым панелям не доступна этому пользователю")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200 OK",
		"message": "заявка по солнечным панелям успешно удалена",
	})
}
