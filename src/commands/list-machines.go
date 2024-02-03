package commands

import (
	"fmt"
	"strings"
)

type MachineStatus string

const (
	MachineStatusStopped  MachineStatus = "stopped"
	MachineStatusStopping MachineStatus = "stopping"
	MachineStatusStarting MachineStatus = "starting"
	MachineStatusStarted  MachineStatus = "started"
	MachineStatusPausing  MachineStatus = "pausing"
	MachineStatusPaused   MachineStatus = "paused"
	MachineStatusResuming MachineStatus = "resuming"
)

type MachineDef struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	Status MachineStatus `json:"status"`
}

type FieldSetter = func(machine *MachineDef, value string)

var noSetter = func(machine *MachineDef, value string) {}
var machineFieldMap = map[string]FieldSetter{
	"uuid":   func(machine *MachineDef, value string) { machine.ID = value },
	"name":   func(machine *MachineDef, value string) { machine.Name = value },
	"status": func(machine *MachineDef, value string) { machine.Status = MachineStatus(value) },
}

func (c *Commander) ListMachines() ([]MachineDef, error) {
	out, err := c.exec([]string{"list"})
	if err != nil {
		return nil, fmt.Errorf("error listing machines: %w", err)
	}

	lines := strings.Split(out, "\n")
	machines := make([]MachineDef, 0, len(lines)-1)

	setters := make([]FieldSetter, 0)

	for lineIx, line := range lines {
		isTitle := lineIx == 0
		var fields []string

		if isTitle {
			fields = strings.Fields(line)
		} else {
			fields = strings.SplitN(line, " ", len(setters))
		}
		var machine MachineDef

		for fieldIx, field := range fields {
			if isTitle {
				fieldTitle := strings.ToLower(field)

				if setter, ok := machineFieldMap[fieldTitle]; ok {
					setters = append(setters, setter)
				} else {
					setters = append(setters, noSetter)
				}

				continue
			}

			// Sets the value of the field in the machine struct based on the detected title
			setters[fieldIx](&machine, strings.TrimSpace(field))
		}

		if !isTitle {
			machines = append(machines, machine)
		}
	}

	return machines, nil
}
