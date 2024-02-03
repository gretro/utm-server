package api

import (
	"context"
	"slices"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gretro/utm_server/src/commands"
	"github.com/gretro/utm_server/src/config"
)

const defaultTimeout = 60 * time.Second

func RegisterMachineStatusRoutes(router gin.IRouter) {
	group := router.Group("/v1/machines/:id/status")
	group.GET("", getMachineStatusV1)
	group.POST("start", startMachineV1)
	group.DELETE("stop", stopMachineV1)
	group.POST("suspend", suspendMachineV1)
}

func getMachineStatusV1(c *gin.Context) {
	commander := commands.New(config.GetAppConfig())
	machine, err := commander.GetMachineByID(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, machine)
}

type StartMachineV1Request struct {
	Disposable       bool `json:"disposable,omitempty"`
	Wait             bool `json:"wait,omitempty"`
	TimeoutInSeconds uint `json:"timeoutInSeconds,omitempty"`
}

func startMachineV1(c *gin.Context) {
	machineID := c.Param("id")

	request := StartMachineV1Request{}
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&request); err != nil {
			c.Error(err)
			return
		}
	}

	timeout := defaultTimeout
	if request.TimeoutInSeconds > 0 {
		timeout = time.Duration(request.TimeoutInSeconds) * time.Second
	}

	commander := commands.New(config.GetAppConfig())
	machine, err := commander.StartMachine(machineID, commands.StartMachineOptions{
		Disposable: request.Disposable,
	})
	if err != nil {
		c.Error(err)
		return
	}

	if request.Wait {
		ctx, cancel := context.WithTimeout(c, timeout)
		defer cancel()

		machine, err = commander.WaitForState(ctx, machineID, commands.MachineStatusStarted)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.JSON(200, machine)
}

func stopMachineV1(c *gin.Context) {
	machineID := c.Param("id")

	rawMethod := c.Query("method")
	method := commands.StopRequest
	if rawMethod != "" && slices.Contains(commands.StopMethods(), rawMethod) {
		method = commands.StopMachineMethod(rawMethod)
	}

	commander := commands.New(config.GetAppConfig())
	machine, err := commander.StopMachine(machineID, commands.StopMachineOptions{
		StopMethod: method,
	})
	if err != nil {
		c.Error(err)
		return
	}

	wait := c.Query("wait") == "true" || c.Query("wait") == "1"
	rawTimeoutInS := c.Query("timeoutInSeconds")

	if wait {
		timeout := defaultTimeout
		if rawTimeoutInS != "" {
			if sec, err := strconv.ParseInt(rawTimeoutInS, 10, 32); err == nil {
				timeout = time.Duration(sec) * time.Second
			}
		}

		ctx, cancel := context.WithTimeout(c, timeout)
		defer cancel()

		machine, err = commander.WaitForState(ctx, machineID, commands.MachineStatusStopped)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.JSON(200, machine)
}

func suspendMachineV1(c *gin.Context) {
	machineID := c.Param("id")
	saveState := c.Query("save-state") == "true" || c.Query("save-state") == "1"

	commander := commands.New(config.GetAppConfig())
	machine, err := commander.SuspendMachine(machineID, commands.SuspendMachineOptions{
		SaveState: saveState,
	})
	if err != nil {
		c.Error(err)
		return
	}

	wait := c.Query("wait") == "true" || c.Query("wait") == "1"
	rawTimeoutInS := c.Query("timeoutInSeconds")

	if wait {
		timeout := defaultTimeout
		if rawTimeoutInS != "" {
			if sec, err := strconv.ParseInt(rawTimeoutInS, 10, 32); err == nil {
				timeout = time.Duration(sec) * time.Second
			}
		}

		ctx, cancel := context.WithTimeout(c, timeout)
		defer cancel()

		machine, err = commander.WaitForState(ctx, machineID, commands.MachineStatusPaused)
		if err != nil {
			c.Error(err)
			return
		}
	}

	c.JSON(200, machine)
}
