package commands

import (
	"fmt"
)

type CloneMachineArgs struct {
	SourceMachineID string
	NewMachineName  string
}

func (c *Commander) CloneMachine(args CloneMachineArgs) (MachineDef, error) {
	out, err := c.exec([]string{"clone", args.SourceMachineID, "--name", args.NewMachineName})
	if err != nil {
		return MachineDef{}, fmt.Errorf("error cloning machine: %w", err)
	}

	// If there is any output, it's likely an error from stderr
	if out != "" {
		return MachineDef{}, fmt.Errorf("error cloning machine: %s", out)
	}

	// There is an issue where we clone a machine and it won't boot at the same time
	// as other instances from the same template. This is because the MAC address is the same.
	// Even changing the MAC address in the plist file doesn't seem to fix it, since UTM likely
	// caches the value of the field.

	machine, err := c.GetMachineByName(args.NewMachineName)
	if err != nil {
		return MachineDef{}, err
	}

	return machine, nil
}
