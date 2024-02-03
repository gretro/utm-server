package commands

import "fmt"

func (c *Commander) DeleteMachine(machineID string) error {
	_, err := c.exec([]string{"delete", machineID})
	if err != nil {
		return fmt.Errorf("error deleting machine: %w", err)
	}

	return nil
}
