package app

var _ ansibleConfig = (*ansibleSetupVMConfig)(nil)

type ansibleSetupVMConfig struct {
	ansibleBaseConfig
}

func newAnsibleSetupVMConfig() *ansibleSetupVMConfig {
	return &ansibleSetupVMConfig{
		ansibleBaseConfig: ansibleBaseConfig{
			SSHPort: "22",
		},
	}
}

func (c *ansibleSetupVMConfig) FillInKeys() error {
	return c.fillInBaseKeys()
}
