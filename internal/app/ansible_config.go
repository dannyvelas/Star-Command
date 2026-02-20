package app

type ansibleConfig interface {
	NodeIP() string
	SSHPort() string
	SSHUser() string
	SSHPrivateKeyPath() string
}
