package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dannyvelas/starcommand/config"
	"github.com/dannyvelas/starcommand/internal/helpers"
	"github.com/goccy/go-yaml"
)

type setupHostEntry struct {
	AnsibleBaseConfig ansibleBaseConfig
	Incus             config.IncusConfig
}

type setupHostHostVars struct {
	AnsibleHost          string `yaml:"ansible_host"`
	AnsiblePort          int    `yaml:"ansible_port"`
	AnsibleSSHPrivateKey string `yaml:"ansible_ssh_private_key_file"`
	AnsibleUser          string `yaml:"ansible_user"`
	IncusStoragePoolName string `yaml:"incus_storage_pool_name"`
	IncusStorageDriver   string `yaml:"incus_storage_driver"`
}

type ansibleSetupHostConfig struct {
	Hosts []setupHostEntry `json:"-" required:"true"`

	// Sensitive
	SMTPUser     string `json:"smtp_user" sensitive:"true" prompt:"SMTP username"`
	SMTPPassword string `json:"smtp_password" sensitive:"true" prompt:"SMTP password"`
}

func newAnsibleSetupHostConfig(c *config.Config) *ansibleSetupHostConfig {
	setupConfig := new(ansibleSetupHostConfig)
	for _, host := range c.Hosts {
		setupConfig.Hosts = append(setupConfig.Hosts, setupHostEntry{
			AnsibleBaseConfig: newAnsibleBaseConfig(host.Name, host.IP, host.SSH),
			Incus:             host.Incus,
		})
	}
	return setupConfig
}

func (c *ansibleSetupHostConfig) generateHostVars() error {
	for _, host := range c.Hosts {
		ansibleUser, err := determineAnsibleUser(host.AnsibleBaseConfig.SSH.User, host.AnsibleBaseConfig.IP, host.AnsibleBaseConfig.SSH.Port, host.AnsibleBaseConfig.SSH.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("error determining ansible user for host %s: %v", host.AnsibleBaseConfig.Name, err)
		}

		expandedPrivateKey, err := helpers.ExpandPath(host.AnsibleBaseConfig.SSH.PrivateKeyPath)
		if err != nil {
			return fmt.Errorf("error expanding private key path for host %s: %v", host.AnsibleBaseConfig.Name, err)
		}

		vars := setupHostHostVars{
			AnsibleHost:          host.AnsibleBaseConfig.IP,
			AnsiblePort:          host.AnsibleBaseConfig.SSH.Port,
			AnsibleSSHPrivateKey: expandedPrivateKey,
			AnsibleUser:          ansibleUser,
			IncusStoragePoolName: host.Incus.StoragePoolName,
			IncusStorageDriver:   host.Incus.StoragePoolDriver,
		}

		dir := filepath.Join(".generated", "host_vars", host.AnsibleBaseConfig.Name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("error creating host_vars dir for %s: %v", host.AnsibleBaseConfig.Name, err)
		}

		data, err := yaml.Marshal(vars)
		if err != nil {
			return fmt.Errorf("error marshaling host vars for %s: %v", host.AnsibleBaseConfig.Name, err)
		}

		if err := os.WriteFile(filepath.Join(dir, "vars.yml"), data, 0o644); err != nil {
			return fmt.Errorf("error writing host vars file for %s: %v", host.AnsibleBaseConfig.Name, err)
		}
	}
	return nil
}
