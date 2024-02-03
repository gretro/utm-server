package commands

import "fmt"

type CloneMachineArgs struct {
	SourceMachineID string
	NewMachineName  string
}

func (c *Commander) CloneMachine(args CloneMachineArgs) (MachineDef, error) {
	_, err := c.exec([]string{"clone", args.SourceMachineID, "--name", args.NewMachineName})
	if err != nil {
		return MachineDef{}, fmt.Errorf("error cloning machine: %w", err)
	}

	machine, err := c.GetMachineByName(args.NewMachineName)
	if err != nil {
		return MachineDef{}, err
	}

	return machine, nil
}
