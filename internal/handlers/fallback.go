package handlers

import (
	"github.com/dannyvelas/homelab/internal/models"
)

func fallbackTargetToStruct(hostAlias, target string) (any, error) {
	switch target {
	case "ssh":
		return models.NewSSHHost(hostAlias), nil
	default:
		return nil, ErrNotFound
	}
}
