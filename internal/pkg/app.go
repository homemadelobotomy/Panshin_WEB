package pkg

import (
	"fmt"

	"lab/internal/app/config"
	"lab/internal/app/handler"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SolarPanelPowerCalculator struct {
	Config  *config.Config
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(c *config.Config, r *gin.Engine, h *handler.Handler) *SolarPanelPowerCalculator {
	return &SolarPanelPowerCalculator{
		Config:  c,
		Router:  r,
		Handler: h,
	}
}

func (a *SolarPanelPowerCalculator) RunApp() {
	logrus.Info("Server start up")

	a.Handler.RegisterHandlers(a.Router)
	serverAddress := fmt.Sprintf("%s:%d", a.Config.ServiceHost, a.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
