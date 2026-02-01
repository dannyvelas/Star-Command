package handlers

type terraformProxmoxConfig struct {
	Node              string `json:"node"`
	Endpoint          string `json:"endpoint"`
	APIToken          string `json:"api_token"`
	SSHAddress        string `json:"ssh_address"`
	SSHPort           string `json:"ssh_port"`
	SSHPrivateKeyPath string `json:"ssh_private_key_path"`
	TerraformVersion  string `json:"terraform_version"`
}

func newTerraformProxmoxConfig() *terraformProxmoxConfig {
	return &terraformProxmoxConfig{
		SSHPort: "22",
	}
}
