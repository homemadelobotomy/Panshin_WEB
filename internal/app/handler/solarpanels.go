package handler

import (
	"lab/internal/app/ds"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetSolarPanels(ctx *gin.Context) {
	var panels []ds.SolarPanel
	var err error

	beginStr := ctx.Query("begin")
	endStr := ctx.Query("end")
	var begin, end int
	if endStr == "" && beginStr == "" {
		panels, err = h.Repository.GetSolarPanels()
		if err != nil {
			logrus.Error(err)
		}

	} else {
		begin, err = strconv.Atoi(beginStr)

		if err != nil {
			logrus.Error("invalid input in begin")
			begin = 0
		}

		end, err = strconv.Atoi(endStr)

		if err != nil {
			logrus.Error("invalid input in end")
			end = 0
		}
		if begin != 0 && end == 0 {
			panels, err = h.Repository.GetSolarPanelsInRange(begin, 100000)
			if err != nil {
				logrus.Error(err)
			}
		} else {
			panels, err = h.Repository.GetSolarPanelsInRange(begin, end)
			if err != nil {
				logrus.Error(err)
			}
		}

	}
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "panel-catalog.html", gin.H{
		"panels": panels,
		"begin":  beginStr,
		"end":    endStr,
		"count":  h.Repository.GetNumberOfPanelsInRequest(1),
	})
}

func (h *Handler) GetSolarPanel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	panel, err := h.Repository.GetSolarPanel(id)
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "panel-details.html", gin.H{
		"panel": panel,
	})
}

func (h *Handler) GetSolarPanelRequest(ctx *gin.Context) {
	request, err := h.Repository.GetSolarPanelRequest(1)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "panel-request.html", gin.H{
		"request":         request,
		"panels_requests": request.Panels,
	})
}

func (h *Handler) AddSolarPanelToRequest(ctx *gin.Context) {
	request_id := ctx.Param("request_id")
	solarpanel_id := ctx.Param("solarpanel_id")

}
