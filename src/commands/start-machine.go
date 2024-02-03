package commands

import "fmt"

type StartMachineOptions struct {
	// Disposable runs VM as a snapshot and does not save changes to disk if set as true.
	Disposable bool
}

func (c *Commander) StartMachine(machineID string, options StartMachineOptions) (Machine, error) {
	cmd := []string{"start", machineID}

	if options.Disposable {
		cmd = append(cmd, "--disposable")
	}

	out, err := c.exec(cmd)
	if err != nil {
		return Machine{}, fmt.Errorf("error starting machine: %w", err)
	}

	// If there is an output, it is likely an error message.
	if out != "" {
		return Machine{}, fmt.Errorf("could not start machine: %s", out)
	}

	machine, err := c.GetMachineByID(machineID)
	if err != nil {
		return Machine{}, err
	}

	return machine, nil
}
