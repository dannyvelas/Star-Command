package handlers

import (
	"fmt"

	"github.com/dannyvelas/homelab/internal/models"
)

func fallbackTargetToStruct(hostAlias, target string) (any, error) {
	switch target {
	case "ssh":
		sshHost, err := models.NewSSHHost(hostAlias)
		if err != nil {
			return nil, fmt.Errorf("error initializing ssh host: %v", err)
		}
		return sshHost, nil
	default:
		return nil, ErrNotFound
	}
}
