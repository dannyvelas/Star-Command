package app

import (
	"fmt"

	"github.com/dannyvelas/homelab/internal/helpers"
)

var _ ansibleConfig = (*ansibleSetupVMConfig)(nil)

type ansibleSetupVMConfig struct {
	// Required
	NodeIP            string `json:"node_ip" required:"true"`
	SSHUser           string `json:"ssh_user" required:"true"`
	SSHPort           string `json:"ssh_port" required:"true"`
	SSHPrivateKeyPath string `json:"ssh_private_key_path" required:"true"`

	// Injected
	AnsibleUser string `json:"ansible_user"`
	AnsiblePort string `json:"ansible_port"`
}

// newAnsibleSetupVMConfig returns a pointer to an ansibleSetupVMConfig struct with some defaults
func newAnsibleSetupVMConfig() *ansibleSetupVMConfig {
	return &ansibleSetupVMConfig{
		SSHPort: "22",
	}
}

func (c *ansibleSetupVMConfig) FillInKeys() error {
	expandedPrivateKeyPath, err := helpers.ExpandPath(c.SSHPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("error expanding path(%s): %v", c.SSHPrivateKeyPath, err)
	}
	c.SSHPrivateKeyPath = expandedPrivateKeyPath

	c.AnsibleUser = c.SSHUser
	c.AnsiblePort = c.SSHPort

	return nil
}

func (c *ansibleSetupVMConfig) GetNodeIP() string {
	return c.NodeIP
}

func (c *ansibleSetupVMConfig) GetSSHUser() string {
	return c.SSHUser
}

func (c *ansibleSetupVMConfig) GetSSHPort() string {
	return c.SSHPort
}

func (c *ansibleSetupVMConfig) GetSSHPrivateKeyPath() string {
	return c.SSHPrivateKeyPath
}
