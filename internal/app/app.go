package app

import (
	"errors"
	"fmt"
	"maps"

	"github.com/dannyvelas/conflux"
	"github.com/dannyvelas/homelab/internal/helpers"
	"github.com/dannyvelas/homelab/internal/models"
	"github.com/go-viper/mapstructure/v2"
)

func GetConfig(hostAlias string, targets []string) (map[string]string, map[string]string, error) {
	configMux := conflux.NewConfigMux(
		conflux.WithYAMLFileReader(helpers.FallbackFile, conflux.WithPath(helpers.GetConfigPath(hostAlias))),
		conflux.WithEnvReader(),
		conflux.WithBitwardenSecretReader(),
	)

	configStructs, err := models.AliasToStruct(hostAlias, targets)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting alias to struct: %v", err)
	}

	allConfigs, allDiagnostics := make(map[string]string), make(map[string]string)
	for _, configStruct := range configStructs {
		diagnostics, err := conflux.Unmarshal(configMux, configStruct)
		if errors.Is(err, conflux.ErrInvalidFields) {
			maps.Copy(allDiagnostics, diagnostics)
			continue
		} else if err != nil {
			return nil, nil, fmt.Errorf("error unmarshalling: %v", err)
		}

		if err := mapstructure.Decode(configStruct, &allConfigs); err != nil {
			return nil, nil, fmt.Errorf("error merging config struct to map: %v", err)
		}
	}

	return allConfigs, allDiagnostics, nil
}

func CheckConfig(hostAlias string, targets []string) (map[string]string, error) {
	configMux := conflux.NewConfigMux(
		conflux.WithYAMLFileReader(helpers.FallbackFile, conflux.WithPath(helpers.GetConfigPath(hostAlias))),
		conflux.WithEnvReader(),
		conflux.WithBitwardenSecretReader(),
	)

	configStructs, err := models.AliasToStruct(hostAlias, targets)
	if err != nil {
		return nil, fmt.Errorf("error getting alias to struct: %v", err)
	}

	allDiagnostics := make(map[string]string)
	for _, configStruct := range configStructs {
		diagnostics, err := conflux.Unmarshal(configMux, configStruct)
		if errors.Is(err, conflux.ErrInvalidFields) {
			maps.Copy(allDiagnostics, diagnostics)
			continue
		} else if err != nil {
			return nil, fmt.Errorf("error unmarshalling: %v", err)
		}
	}

	return allDiagnostics, nil
}
