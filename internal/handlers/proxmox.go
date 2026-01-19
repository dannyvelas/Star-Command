package handlers

import "github.com/dannyvelas/homelab/internal/models"

type ProxmoxHandler struct {
	hostAlias string
}

func NewProxmoxHandler() ProxmoxHandler {
	return ProxmoxHandler{
		hostAlias: "proxmox",
	}
}

func (h ProxmoxHandler) TargetToStruct(target string) (any, error) {
	switch target {
	case "ansible":
		return models.NewAnsibleProxmoxConfig(), nil
	default:
		return fallbackTargetToStruct(h.hostAlias, target)
	}
}
