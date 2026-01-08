package config

import (
	"fmt"

	"github.com/dannyvelas/homelab/internal/helpers"
)

const (
	statusMissing = "missing"
	statusLoaded  = "loaded"
)

type config interface {
	// Validate returns an map of validation results where each element corresponds to a key in the config
	// the second return value will be false if at least one key was invalid. otherwise, it will be true
	Validate(map[string]string, bool) (map[string]string, bool)
}

type fillableConfig interface {
	// FillInKeys takes the keys that are required and uses them to fill out remaining config fields
	FillInKeys() error
}

func validateStruct(v any) (map[string]string, bool, error) {
	results := make(map[string]string)
	ok := true

	tagToFieldMap, err := helpers.GetTagToFieldMap(v, "bw", "json")
	if err != nil {
		return nil, false, fmt.Errorf("error getting tag to field map: %v", err)
	}

	for tag, field := range tagToFieldMap {
		if _, ok := field.Type.Tag.Lookup("required"); !ok {
			continue
		}

		if field.Value.IsZero() {
			results[tag] = statusMissing
			ok = false
		} else {
			results[tag] = statusLoaded
		}
	}

	config, ok := v.(config)
	if !ok {
		return results, ok, nil
	}

	results, ok = config.Validate(results, ok)
	return results, ok, nil
}
