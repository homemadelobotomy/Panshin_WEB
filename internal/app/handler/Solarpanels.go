package handler

import (
	"errors"
	dto "lab/internal/app/DTO"
	"lab/internal/app/ds"
	"lab/internal/app/service"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

func (h *Handler) RegisterSolarPanelHandlers(router *gin.Engine) {
	panelsGroup := router.Group("/api/panels")
	{
		panelsGroup.GET("", h.GetSolarPanels)
		panelsGroup.GET("/:id", h.GetOneSolarPanel)
		panelsGroup.POST("", h.AddNewSolarPanel)
		panelsGroup.PUT("/:id", h.ChangeSolarPanel)
		panelsGroup.DELETE("/:id", h.DeleteSolarPanel)
		panelsGroup.POST("/:id", h.AddSolarPanelToRequest)
		panelsGroup.POST("/:id/image", h.AddImageToSolarPanel)
	}
}

func (h *Handler) GetSolarPanels(ctx *gin.Context) {
	var (
		startValue float64
		endValue   float64
		err        error
	)
	startValueStr := ctx.Query("start_value")
	endValueStr := ctx.Query("end_value")

	if startValueStr != "" {
		startValue, err = strconv.ParseFloat(startValueStr, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено некорректное значение start_value")
		}
	} else {
		startValue = 0
	}

	if endValueStr != "" {
		endValue, err = strconv.ParseFloat(endValueStr, 64)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено некорректное значение end_value")
		}
	} else {
		endValue = 0
	}

	response, err := h.Service.GetSolarPanels(startValue, endValue)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявки по солнечным панелям с таким id не найден")
		} else if errors.Is(err, service.ErrNoRecords) {
			h.errorHandler(ctx, http.StatusNotFound,
				"по заданным параметрам не найдено солнечных панелей")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"введено недопсутимое значение мощности для фильтрации")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) GetOneSolarPanel(ctx *gin.Context) {
	solarPanelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}
	response, err := h.Service.GetSolarPanel(uint(solarPanelId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечной панели с таким id не найдено")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) AddNewSolarPanel(ctx *gin.Context) {
	var solarPanel dto.AddSolarPanel
	if err := ctx.BindJSON(&solarPanel); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса")
		return
	}
	response, err := h.Service.AddNewSolarPanel(solarPanel)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	ctx.JSON(http.StatusCreated, response)
}

func (h *Handler) ChangeSolarPanel(ctx *gin.Context) {
	var solarPanel dto.ChangeSolarPanel
	solarPanelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}
	if err := ctx.BindJSON(&solarPanel); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest,
			"введен неправильный формат тела запроса: "+err.Error())
		return
	}
	response, err := h.Service.ChangeSolarPanel(uint(solarPanelId), solarPanel)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечной панели с таким id не найдено")
		} else if errors.Is(err, service.ErrBadRequest) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"числовые значения должны быть положительными")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSolarPanel(ctx *gin.Context) {
	solarPanelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}
	err = h.Service.DeleteSolarPanel(uint(solarPanelId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечной панели с таким id не найдено")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"status":  "200 OK",
		"message": "солнечная панель успешно удалена",
	})
}

func (h *Handler) AddSolarPanelToRequest(ctx *gin.Context) {
	solarPanelId, err := strconv.Atoi(ctx.Param("id"))
	var pgErr *pgconn.PgError
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}
	err = h.Service.AddSolarPanelToRequest(uint(solarPanelId), ds.GetUser().GetId())
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"заявка по слонечным панелям не найдена")
		} else if errors.Is(err, service.ErrNoRecords) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечной панели с таким id не найдено")
		} else if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			h.errorHandler(ctx, http.StatusBadRequest,
				"такая солнечная панель уже добавлена в заявку")
		} else if errors.Is(err, service.ErrSolarPanelDeleted) {
			h.errorHandler(ctx, http.StatusBadRequest,
				"эта солнечная панель удалена")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "201 Created",
		"message": "солнечная панель добавлена в заявку",
	})
}

func (h *Handler) AddImageToSolarPanel(ctx *gin.Context) {
	solarPanelId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusNotFound,
			"введен некорректный id солнечной панели")
		return
	}
	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, "файл не найден")
		return
	}

	filename := uuid.New().String() + filepath.Ext(file.Filename)

	response, err := h.Service.AddImageToSolarPanel(uint(solarPanelId), file, filename)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.errorHandler(ctx, http.StatusNotFound,
				"солнечной панели с таким id не найдено")
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
		}
		return
	}
	ctx.JSON(http.StatusOK, response)

}

// func (h *Handler) GetSolarPanels(ctx *gin.Context) {
// 	var panels []ds.SolarPanel
// 	var err error

// 	beginStr := ctx.Query("begin")
// 	endStr := ctx.Query("end")
// 	var begin, end int
// 	if endStr == "" && beginStr == "" {
// 		panels, err = h.Repository.GetSolarPanels()
// 		if err != nil {
// 			logrus.Error(err)
// 		}

// 	} else {
// 		begin, err = strconv.Atoi(beginStr)
// 		if err != nil {
// 			logrus.Error("invalid input in begin")
// 			begin = 0
// 		}
// 		end, err = strconv.Atoi(endStr)
// 		if err != nil {
// 			logrus.Error("invalid input in end")
// 			end = 0
// 		}
// 		if begin != 0 && end == 0 {
// 			panels, err = h.Repository.GetSolarPanelsInRange(begin, 100000)
// 			if err != nil {
// 				logrus.Error(err)
// 			}
// 		} else {
// 			panels, err = h.Repository.GetSolarPanelsInRange(begin, end)
// 			if err != nil {
// 				logrus.Error(err)
// 			}
// 		}
// 	}
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	request_id, err := h.Repository.GetSolarPanelRequestID()
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	count, err := h.Repository.GetNumberOfPanelsInRequest()
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "panel-catalog.html", gin.H{
// 		"panels":     panels,
// 		"begin":      beginStr,
// 		"end":        endStr,
// 		"count":      count,
// 		"request_id": request_id,
// 	})
// }

// func (h *Handler) GetSolarPanel(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	panel, err := h.Repository.GetSolarPanel(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	ctx.HTML(http.StatusOK, "panel-details.html", gin.H{
// 		"panel": panel,
// 	})
// }

// func (h *Handler) GetSolarPanelRequest(ctx *gin.Context) {
// 	request_id, err := strconv.Atoi(ctx.Param("id"))
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	request, err := h.Repository.GetSolarPanelRequest(request_id)
// 	if err != nil {
// 		ctx.Redirect(http.StatusSeeOther, "/panels")
// 	}

// 	ctx.HTML(http.StatusOK, "panel-request.html", gin.H{
// 		"request":         request,
// 		"panels_requests": request.Panels,
// 	})
// }

// func (h *Handler) AddSolarPanelToRequest(ctx *gin.Context) {

// 	solarpanel_id, err := strconv.Atoi(ctx.Param("solarpanel_id"))
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	h.Repository.AddSolarPanelToRequest(solarpanel_id)

// 	ctx.Redirect(http.StatusSeeOther, "/panels")
// }

// func (h *Handler) DeleteSolarPanelRequest(ctx *gin.Context) {
// 	request_id, err := strconv.Atoi(ctx.Param("request_id"))
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	err = h.Repository.DeleteSolarPanelRequest(request_id)
// 	if err != nil {
// 		h.errorHandler(ctx, http.StatusInternalServerError, err)
// 		return
// 	}
// 	ctx.Redirect(http.StatusSeeOther, "/panels")
// }
