package commands

import "fmt"

type StopMachineMethod string

const (
	// StopForce stops by sending a power off event (default)
	StopForce StopMachineMethod = "force"
	// StopKill stops by killing the VM process
	StopKill StopMachineMethod = "kill"
	// StopRequest stops by sending a request to the VM
	StopRequest StopMachineMethod = "request"
)

func StopMethods() []string {
	return []string{
		string(StopForce),
		string(StopKill),
		string(StopRequest),
	}
}

type StopMachineOptions struct {
	// StopMethod is the method used to stop the machine
	StopMethod StopMachineMethod
}

func (c *Commander) StopMachine(machineID string, options StopMachineOptions) (Machine, error) {
	cmd := []string{"stop", machineID}

	switch options.StopMethod {
	case StopForce:
		cmd = append(cmd, "--force")
	case StopKill:
		cmd = append(cmd, "--kill")
	case StopRequest:
		cmd = append(cmd, "--request")
	}

	_, err := c.exec(cmd)
	if err != nil {
		return Machine{}, fmt.Errorf("error stopping machine: %w", err)
	}

	machine, err := c.GetMachineByID(machineID)
	if err != nil {
		return Machine{}, err
	}

	return machine, nil
}
