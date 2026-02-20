package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/crypto/ssh"
)

var _ Handler = AnsibleProxmoxHandler{}

type AnsibleProxmoxHandler struct{}

func NewAnsibleProxmoxHandler() AnsibleProxmoxHandler {
	return AnsibleProxmoxHandler{}
}

func (h AnsibleProxmoxHandler) GetConfig(_ string) any {
	return newAnsibleProxmoxConfig()
}

func (h AnsibleProxmoxHandler) Execute(_ context.Context, config any, hostAlias string) (map[string]string, error) {
	diagnostics := make(map[string]string)

	ansibleProxmoxConfig, ok := config.(*ansibleProxmoxConfig)
	if !ok {
		return diagnostics, fmt.Errorf("internal type error converting config to ansible proxmox config. found: %T", config)
	}

	if err := h.runAnsiblePlaybook(ansibleProxmoxConfig); err != nil {
		return diagnostics, fmt.Errorf("error running ansible playbook: %v", err)
	}

	return diagnostics, nil
}

func (h AnsibleProxmoxHandler) runAnsiblePlaybook(config *ansibleProxmoxConfig) error {
	proxmoxAddr := fmt.Sprintf("%s:%s", config.NodeIP, config.SSHPort)
	client, sshErr := h.getSSHClient(config.SSHUser, proxmoxAddr, config.SSHPrivateKeyPath)
	if sshErr != nil && !errors.Is(sshErr, errConnectingSSH) {
		return fmt.Errorf("error checking if ssh is accessible to proxmox host: %v", sshErr)
	} else if sshErr == nil {
		_ = client.Close()
	}

	tmpFile, err := os.CreateTemp("", "labctl-vars-*.json")
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	defer func() { _ = os.Remove(tmpFile.Name()) }()

	if err := json.NewEncoder(tmpFile).Encode(config); err != nil {
		return fmt.Errorf("error writing config to tmp file: %v", err)
	}

	args := []string{"-i", "ansible/inventory.ini", "ansible/setup-proxmox.yml", "-e", "@" + tmpFile.Name()}
	if errors.Is(sshErr, errConnectingSSH) {
		args = append(args, "-u", "root")
	}

	cmd := exec.Command("ansible-playbook", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running ansible proxmox command: %v", err)
	}

	return nil
}

func (h AnsibleProxmoxHandler) getSSHClient(user, addr, privateKeyPath string) (*ssh.Client, error) {
	key, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read private key: %v", err)
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("unable to parse private key: %v", err)
	}

	sshClientConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		Timeout:         3 * time.Second,             // -o ConnectTimeout=3
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // StrictHostKeyChecking=no
	}

	client, err := ssh.Dial("tcp", addr, sshClientConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errConnectingSSH, err)
	}

	return client, nil
}
