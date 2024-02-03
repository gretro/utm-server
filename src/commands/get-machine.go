package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/gretro/utm_server/src/system"
)

type IPAddressType string

const (
	IPAddressTypeUnknown IPAddressType = "unknown"
	IPAddressTypeIPv4    IPAddressType = "ipv4"
	IPAddressTypeIPv6    IPAddressType = "ipv6"
)

var (
	ipv4Regex = regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`)
	ipv6Regex = regexp.MustCompile(`[a-fA-F0-9:]+`)
)

type IPAddress struct {
	Type IPAddressType `json:"type"`
	IP   string        `json:"ip"`
}

type Machine struct {
	MachineDef
	IPAddresses []IPAddress `json:"ipAddresses"`
}

func (c *Commander) GetMachineByID(machineID string) (Machine, error) {
	machines, err := c.ListMachines()
	if err != nil {
		return Machine{}, err
	}

	for _, machine := range machines {
		if machine.ID == machineID {
			ipAddresses := make([]IPAddress, 0)

			if machine.Status == MachineStatusStarted {
				ipAddresses, _ = c.GetMachineIPAddresses(machineID)
			}

			return Machine{
				MachineDef:  machine,
				IPAddresses: ipAddresses,
			}, nil
		}
	}

	return Machine{}, ErrMachineNotFound
}

func (c *Commander) GetMachineIPAddresses(machineID string) ([]IPAddress, error) {
	out, err := c.exec([]string{"ip-address", machineID})
	if err != nil {
		c.l.Debug("error getting IP addresses", "machine", machineID, system.ErrorLabel, err)
		return nil, fmt.Errorf("error getting IP addresses: %w", err)
	}

	lines := strings.Split(out, "\n")
	ips := make([]IPAddress, 0, len(lines))

	for _, line := range lines {
		addr := strings.TrimSpace(line)

		addrType := IPAddressTypeUnknown
		if ipv4Regex.MatchString(addr) {
			addrType = IPAddressTypeIPv4
		} else if ipv6Regex.MatchString(addr) {
			addrType = IPAddressTypeIPv6
		}

		ips = append(ips, IPAddress{
			Type: addrType,
			IP:   addr,
		})
	}

	return ips, nil
}

func (c *Commander) GetMachineByName(machineName string) (MachineDef, error) {
	machines, err := c.ListMachines()
	if err != nil {
		return MachineDef{}, err
	}

	for _, machine := range machines {
		if machine.Name == machineName {
			return machine, nil
		}
	}

	return MachineDef{}, ErrMachineNotFound
}
