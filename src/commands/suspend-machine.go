package commands

import "fmt"

type SuspendMachineOptions struct {
	// SaveState saves the VM state to disk after suspending
	SaveState bool
}

func (c *Commander) SuspendMachine(machineID string, options SuspendMachineOptions) (Machine, error) {
	cmd := []string{"suspend", machineID}

	if options.SaveState {
		cmd = append(cmd, "--save-state")
	}

	_, err := c.exec(cmd)
	if err != nil {
		return Machine{}, fmt.Errorf("error suspending machine: %w", err)
	}

	machine, err := c.GetMachineByID(machineID)
	if err != nil {
		return Machine{}, err
	}

	return machine, nil
}
