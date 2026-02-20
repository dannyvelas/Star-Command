package app

var _ ansibleConfig = (*ansibleSetupHostConfig)(nil)

type ansibleSetupHostConfig struct {
	ansibleBaseConfig

	// Required
	IncusStoragePoolName string `json:"incus_storage_pool_name" required:"true"`
	IncusStorageDriver   string `json:"incus_storage_driver" required:"true"`
}

func newAnsibleSetupHostConfig() *ansibleSetupHostConfig {
	return &ansibleSetupHostConfig{
		ansibleBaseConfig: ansibleBaseConfig{
			SSHPort: "22",
		},
	}
}

func (c *ansibleSetupHostConfig) FillInKeys() error {
	return c.fillInBaseKeys()
}
