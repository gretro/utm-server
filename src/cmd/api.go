package main

import (
	"github.com/gretro/utm_server/src/config"
	"github.com/gretro/utm_server/src/detectos"
	"github.com/gretro/utm_server/src/libs"
	"github.com/gretro/utm_server/src/system"
)

func main() {
	detectos.AssertDarwin()

	appConfig := config.Bootstrap(system.GetComponentLogger("config"))

	libs.BootstrapGin(appConfig)
}
