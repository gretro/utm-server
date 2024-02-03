package libs

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gretro/utm_server/src/api"
	"github.com/gretro/utm_server/src/config"
	"github.com/gretro/utm_server/src/system"
)

func BootstrapGin(appConfig *config.AppConfig) {
	l := system.GetComponentLogger("webserver")

	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// Register routes
	api.RegisterMachineRoutes(engine)
	api.RegisterMachineStatusRoutes(engine)

	// Start the server
	addr := fmt.Sprintf("%s:%d", appConfig.HTTPHost, appConfig.HTTPPort)

	l.Info("Starting server", "addr", addr)
	err := engine.Run(addr)

	if err != nil && !errors.Is(http.ErrServerClosed, err) {
		l.Error("Failed to start server", system.ErrorLabel, err)
		panic("Failed to start server")
	}
}
