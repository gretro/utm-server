package commands

import (
	"context"
	"fmt"
	"time"
)

const timeBetweenPolls = 5 * time.Second

func (c *Commander) WaitForState(ctx context.Context, machineID string, state MachineStatus) (Machine, error) {
	for {
		machine, err := c.GetMachineByID(machineID)
		if err != nil {
			return Machine{}, err
		}

		if machine.Status == state {
			return machine, nil
		}

		c.l.Debug("waiting for machine to reach state", "machine", machineID, "state", state, "current_state", machine.Status)

		select {
		case <-ctx.Done():
			return Machine{}, fmt.Errorf("timed out waiting for machine to reach state %s: %w", state, ctx.Err())
		case <-time.After(timeBetweenPolls):
			continue
		}
	}
}
