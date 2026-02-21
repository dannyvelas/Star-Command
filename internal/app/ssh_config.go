package app

type sshConfig struct {
	Alias         string `json:"alias" required:"true"`
	HostName      string `json:"host_name" required:"true"`
	User          string `json:"ssh_user" required:"true"`
	PublicKeyPath string `json:"ssh_public_key_path" required:"true"`
	Port          string `json:"ssh_port" required:"true"`
}

func newSSHHost(hostAlias string) *sshConfig {
	return &sshConfig{Alias: hostAlias}
}
