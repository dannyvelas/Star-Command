package app

import "github.com/dannyvelas/starcommand/config"

type ansibleBaseConfig struct {
	Name string
	IP   string
	SSH  config.SSHConfig
}

func newAnsibleBaseConfig(name, ip string, ssh config.SSHConfig) ansibleBaseConfig {
	return ansibleBaseConfig{
		Name: name,
		IP:   ip,
		SSH:  ssh,
	}
}
