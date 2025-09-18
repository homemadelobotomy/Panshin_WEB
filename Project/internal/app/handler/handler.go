package handler

import (
	"Project/internal/app/repository"
	"net/http"
	"strconv"

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
func (h *Handler) GetSolarpanels(ctx *gin.Context) {
	var panels []repository.Solarpanel
	var err error
	beginStr := ctx.Query("begin")
	endStr := ctx.Query("end")
	var begin, end int

	if endStr == "" && beginStr == "" {
		panels, err = h.Repository.GetSolarpanels()
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
	panelsInCart, err := h.Repository.GetBid(1)
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "panel-catalog.html", gin.H{
		"panels":      panels,
		"begin":       beginStr,
		"end":         endStr,
		"num_in_cart": len(panelsInCart),
	})
}

func (h *Handler) GetSolarpanel(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}
	panel, err := h.Repository.GetSolarpanel(id)
	if err != nil {
		logrus.Error(err)
	}
	ctx.HTML(http.StatusOK, "panel-details.html", gin.H{
		"panel": panel,
	})
}

func (h *Handler) GetBid(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
	}

	bid, err := h.Repository.GetBid(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "panel-bid.html", gin.H{
		"bid": bid,
	})
}
