package main

import (
	"fmt"
	"os"

	"github.com/gretro/utm_server/src/config"
	"github.com/gretro/utm_server/src/detectos"
	"github.com/gretro/utm_server/src/libs"
	"github.com/gretro/utm_server/src/system"
	"github.com/gretro/utm_server/src/version"
	"github.com/urfave/cli/v2"
)

func main() {
	detectos.AssertDarwin()

	l := system.SystemLogger()
	appConfig := config.NewAppConfig()

	app := &cli.App{
		Name:    "UTM Server",
		Usage:   "Run UTM as a server",
		Version: version.AppVersion(),
		Flags: []cli.Flag{
			cli.HelpFlag,
		},
		Commands: []*cli.Command{
			{
				Name:  "serve",
				Usage: "Starts the UTM server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "utm-app",
						Usage:       "Path to UTM.app",
						Value:       "/Applications/UTM.app",
						DefaultText: "/Applications/UTM.app",
						EnvVars:     []string{"UTM_PATH"},
						Category:    "UTM",
						Destination: &appConfig.UTMPath,
					},
					&cli.StringFlag{
						Name:        "host",
						Usage:       "Host to bind to",
						Value:       "127.0.0.1",
						DefaultText: "127.0.0.1",
						Category:    "HTTP Server",
						Destination: &appConfig.HTTPHost,
					},
					&cli.UintFlag{
						Name:        "port",
						Usage:       "Port to bind to",
						Value:       8788,
						DefaultText: "8788",
						Category:    "HTTP Server",
						Action: func(ctx *cli.Context, port uint) error {
							if port > 65535 {
								return fmt.Errorf("port number %d is invalid", port)
							}

							appConfig.HTTPPort = uint16(port)
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					fmt.Println("UTM Server " + c.App.Version)
					libs.BootstrapGin(appConfig)

					return nil
				},
			},
		},
		DefaultCommand: "serve",
	}

	if err := app.Run(os.Args); err != nil {
		l.Error("Application error", system.ErrorLabel, err)
		os.Exit(1)
	}
}
