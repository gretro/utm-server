package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gretro/utm_server/src/commands"
	"github.com/gretro/utm_server/src/config"
)

func RegisterMachineRoutes(router gin.IRouter) {
	group := router.Group("/v1/machines")
	group.GET("", listMachinesV1)
	group.GET(":id", getMachineV1)
	group.POST("", createMachineV1)
	group.DELETE(":id", deleteMachineV1)
}

type ListMachinesV1Response struct {
	Results []commands.MachineDef `json:"results"`
	Total   int                   `json:"total"`
}

func listMachinesV1(c *gin.Context) {
	commander := commands.New(config.GetAppConfig())
	machines, err := commander.ListMachines()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, ListMachinesV1Response{
		Results: machines,
		Total:   len(machines),
	})
}

func getMachineV1(c *gin.Context) {
	commander := commands.New(config.GetAppConfig())
	machine, err := commander.GetMachineByID(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(200, machine)
}

type CreateMachineV1Request struct {
	Name              string `json:"name,omitempty" binding:"required"`
	TemplateMachineID string `json:"templateMachineId,omitempty" binding:"required"`
}

func createMachineV1(c *gin.Context) {
	req := CreateMachineV1Request{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	commander := commands.New(config.GetAppConfig())
	machine, err := commander.CloneMachine(commands.CloneMachineArgs{
		SourceMachineID: req.TemplateMachineID,
		NewMachineName:  req.Name,
	})

	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(201, machine)
}

func deleteMachineV1(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.Status(http.StatusBadRequest)
		return
	}

	commander := commands.New(config.GetAppConfig())
	err := commander.DeleteMachine(c.Param("id"))
	if err != nil {
		c.Error(err)
		return
	}

	c.Status(204)
}
