package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) RegisterSolarPanelHandlers(router *gin.Engine) {
	router.GET("/panels", h.GetSolarPanels)
	//router.GET("/panel/:id", h.GetSolarPanel)
}

func (h *Handler) GetSolarPanels(ctx *gin.Context) {
	response, err := h.Service.GetSolarPanels()
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err.Error())
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
