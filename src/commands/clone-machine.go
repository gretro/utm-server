package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gretro/utm_server/src/utils"
	"golang.org/x/exp/rand"
)

type CloneMachineArgs struct {
	SourceMachineID string
	NewMachineName  string
}

func (c *Commander) CloneMachine(args CloneMachineArgs) (MachineDef, error) {
	// Cloning the machine via the command line brings up issues with duplicate MAC addresses.
	// Manually copying the VM configuration and fixing issue up.
	c.l.Info("Validation", "sourceMachineID", args.SourceMachineID, "newMachineName", args.NewMachineName)

	templateMachine, err := c.GetMachineByID(args.SourceMachineID)
	if err != nil {
		return MachineDef{}, err
	}

	if templateMachine.Status != MachineStatusStopped {
		return MachineDef{}, fmt.Errorf("source machine must be stopped")
	}

	vmPath, err := utils.ResolvePath(c.vmPath)
	if err != nil {
		return MachineDef{}, fmt.Errorf("failed to resolve VM path: %w", err)
	}

	templateMachinePath := path.Join(vmPath, fmt.Sprintf("%s.utm", templateMachine.Name))
	destinationPath := path.Join(vmPath, fmt.Sprintf("%s.utm", args.NewMachineName))

	_, err = os.Stat(templateMachinePath)
	if err != nil && os.IsNotExist(err) {
		return MachineDef{}, fmt.Errorf("unable to find template machine configuration: %w", err)
	}

	_, err = os.Stat(destinationPath)
	if err == nil {
		return MachineDef{}, fmt.Errorf("destination machine already exists")
	}

	c.l.Info("Cloning source VM", "sourceMachineID", args.SourceMachineID, "newMachineName", args.NewMachineName)
	if err = cloneVMConfig(templateMachinePath, destinationPath); err != nil {
		return MachineDef{}, err
	}

	c.l.Info("Replacing configuration for the new VM", "sourceMachineID", args.SourceMachineID, "newMachineName", args.NewMachineName)
	newMachineID, err := c.replaceConfig(destinationPath, args.NewMachineName)
	if err != nil {
		return MachineDef{}, fmt.Errorf("failed to replace configuration: %w", err)
	}

	c.l.Info("Adding VM to UTM", "sourceMachineID", args.SourceMachineID, "newMachineName", args.NewMachineName)
	if err = c.addVMToUTM(destinationPath); err != nil {
		return MachineDef{}, err
	}

	newMachine, err := c.GetMachineByID(newMachineID)
	if err != nil {
		return MachineDef{}, err
	}

	return newMachine.MachineDef, nil
}

func cloneVMConfig(src, dst string) error {
	cmd := exec.Command("cp", "-r", src, dst)
	out, err := cmd.CombinedOutput()
	if err != nil || len(out) > 0 {
		return fmt.Errorf("failed to clone virtual machine: %s", string(out))
	}

	return nil
}

func generateRandomMAC() string {
	rSource := rand.NewSource(uint64(time.Now().Unix()))
	r := rand.New(rSource)

	buf := make([]byte, 6)
	_, err := r.Read(buf)
	if err != nil {
		panic(err)
	}

	// Set the local bit, necessary for locally administered MAC addresses
	buf[0] |= 2

	mac := fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return strings.ToUpper(mac)
}

func (c *Commander) replaceConfig(vmPath string, machineName string) (string, error) {
	newMachineID := strings.ToUpper(uuid.NewString())

	configFilePath := path.Join(vmPath, "config.plist")

	config, err := os.ReadFile(configFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read cloned machine configuration: %w", err)
	}

	pipedCmds := []*exec.Cmd{
		exec.Command("plutil", "-replace", "Information.Name", "-string", machineName, "-o", "-", "-"),
		exec.Command("plutil", "-replace", "Information.UUID", "-string", newMachineID, "-o", "-", "-"),
	}

	networkCount, err := c.getNetworkCount(configFilePath)
	if err != nil {
		c.l.Warn("Failed to get network count. There may not be any network configurations", "error", err)
		networkCount = 0
	}

	for i := 0; i < networkCount; i++ {
		newMAC := generateRandomMAC()
		pipedCmds = append(
			pipedCmds,
			exec.Command("plutil", "-replace", fmt.Sprintf("Network.%d.MacAddress", i), "-string", newMAC, "-o", "-", "-"),
		)
	}

	if len(pipedCmds) > 0 {
		pipedCmds[0].Stdin = strings.NewReader(string(config))
	}

	transformed, err := c.pipedExec(pipedCmds...)
	if err != nil {
		return "", wrapCmdError(c.l, err)
	}

	err = os.WriteFile(configFilePath, transformed, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write the replacement configuration: %w", err)
	}

	return newMachineID, nil
}

func (c *Commander) getNetworkCount(configFilePath string) (int, error) {
	cmd := exec.Command("plutil", "-extract", "Network", "raw", "-o", "-", configFilePath)
	out, err := cmd.Output()
	if err != nil {
		return 0, wrapCmdError(c.l, err)
	}

	strOut := strings.TrimSpace(string(out))

	count, err := strconv.ParseInt(strOut, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse network configuration count: %w", err)
	}

	return int(count), nil
}

func (c *Commander) addVMToUTM(vmPath string) error {
	cmd := exec.Command("open", vmPath)
	err := cmd.Run()
	if err != nil {
		return wrapCmdError(c.l, err)
	}

	return nil
}
