package handlers

import (
	"fmt"
	"maps"
	"os"
	"slices"

	"github.com/dannyvelas/conflux"
	"github.com/spf13/afero"
)

type Handler interface {
	GetConfig() (map[string]string, map[string]string, error)
	CheckConfig() (map[string]string, error)
	SetFile() ([]string, error)

	useFS(afero.Fs)
	useHomeDir(string)
	getHomeDir() string
	targetToStruct(target string) (any, error)
}

type WritableFile interface {
	Name() string
	Resource() string
	ContentAlreadyExists(fs afero.Fs, homeDir string) (bool, error)
	SetFile(fs afero.Fs, homeDir string) error
}

type HandlerConstructor func(configMux *conflux.ConfigMux, targets []string) Handler

var handlerMap = map[string]HandlerConstructor{
	"proxmox": func(configMux *conflux.ConfigMux, targets []string) Handler {
		return NewProxmoxHandler(configMux, targets)
	},
}

func New(configMux *conflux.ConfigMux, hostAlias string, targets []string, opts ...func(Handler)) (Handler, error) {
	handlerFn, ok := handlerMap[hostAlias]
	if !ok {
		return nil, fmt.Errorf("error: %w: unsupported host(%s)", ErrInvalidArgs, hostAlias)
	}

	handler := handlerFn(configMux, targets)
	for _, opt := range opts {
		opt(handler)
	}

	if handler.getHomeDir() != "" {
		return handler, nil
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error getting user home directory: %v", err)
	}
	handler.useHomeDir(homeDir)

	return handler, nil
}

func GetSupportedHostAliases() []string {
	return slices.Collect(maps.Keys(handlerMap))
}

func WithFS(fs afero.Fs) func(Handler) {
	return func(handler Handler) {
		handler.useFS(fs)
	}
}

func WithHomeDir(homeDir string) func(Handler) {
	return func(handler Handler) {
		handler.useHomeDir(homeDir)
	}
}
