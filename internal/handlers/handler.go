package handlers

import (
	"fmt"
	"maps"
	"slices"

	"github.com/dannyvelas/conflux"
)

type Handler interface {
	GetConfig() (map[string]string, map[string]string, error)
	CheckConfig() (map[string]string, error)
	SetFile() ([]string, error)
	targetToStruct(target string) (any, error)
}

type WritableFile interface {
	SetFile() error
}

type HandlerConstructor func(configMux *conflux.ConfigMux, targets []string) Handler

var handlerMap = map[string]HandlerConstructor{
	"proxmox": func(configMux *conflux.ConfigMux, targets []string) Handler {
		return NewProxmoxHandler(configMux, targets)
	},
}

func New(configMux *conflux.ConfigMux, hostAlias string, targets []string) (Handler, error) {
	handlerFn, ok := handlerMap[hostAlias]
	if !ok {
		return nil, fmt.Errorf("error: %w: unsupported host(%s)", ErrInvalidArgs, hostAlias)
	}

	return handlerFn(configMux, targets), nil
}

func GetSupportedHostAliases() []string {
	return slices.Collect(maps.Keys(handlerMap))
}
