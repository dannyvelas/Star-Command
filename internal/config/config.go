package config

import (
	"errors"
	"fmt"

	"github.com/dannyvelas/homelab/internal/helpers"
)

const (
	StatusMissing = "missing"
	StatusLoaded  = "loaded"
)

type validatable interface {
	// Validate receives a diagnostic map where each element corresponds to a key in the config
	// the second return value will be false if at least one key was invalid. otherwise, it will be true
	Validate(map[string]string) bool
}

type fillable interface {
	// FillInKeys takes the keys that are required and uses them to fill out remaining config fields
	FillInKeys() error
}

func validateStruct(v any) (map[string]string, error) {
	diagnosticMap := make(map[string]string)
	valid := true

	tagToFieldMap, err := helpers.GetTagToFieldMap(v, "labctl", "json")
	if err != nil {
		return nil, fmt.Errorf("error getting tag to field map: %v", err)
	}

	for tag, field := range tagToFieldMap {
		if _, ok := field.Type.Tag.Lookup("required"); !ok {
			continue
		}

		if field.Value.IsZero() {
			diagnosticMap[tag] = StatusMissing
			valid = false
		} else {
			diagnosticMap[tag] = StatusLoaded
		}
	}

	if config, ok := v.(validatable); ok {
		valid = valid && config.Validate(diagnosticMap)
	}

	if !valid {
		return diagnosticMap, ErrInvalidFields
	}

	return diagnosticMap, nil
}

func UnmarshalIntoStruct(r Reader, target any) (map[string]string, error) {
	readResult, err := r.read()
	if err != nil && !errors.Is(err, ErrInvalidFields) {
		return nil, fmt.Errorf("error reading: %v", err)
	} // if errors.Is(err, ErrInvalidFields) we want to continue because we can want to show a report of all missing diagnostics

	if err := helpers.FromMap(readResult.getConfigMap(), target); err != nil {
		return nil, fmt.Errorf("error converting map into target: %v", err)
	}

	readDiagnosticMap := getDiagnosticMap(readResult)

	targetDiagnosticMap, err := validateStruct(target)
	if err != nil && !errors.Is(err, ErrInvalidFields) {
		return nil, fmt.Errorf("error unmarhsalling into config: %v", err)
	}

	mergedDiagnostics := helpers.MergeMaps(readDiagnosticMap, targetDiagnosticMap)
	if errors.Is(err, ErrInvalidFields) {
		return nil, fmt.Errorf("invalid or missing fields:\n%s", diagnosticMapToTable(mergedDiagnostics))
	}

	if fillableTarget, ok := target.(fillable); ok {
		if err := fillableTarget.FillInKeys(); err != nil {
			return nil, fmt.Errorf("error filling in fields: %v", err)
		}
	}

	return mergedDiagnostics, err
}

func unmarshalIntoMap(r Reader, target *map[string]string) (map[string]string, error) {
	readResult, err := r.read()
	if err != nil && !errors.Is(err, ErrInvalidFields) {
		return nil, fmt.Errorf("error reading: %v", err)
	}

	diagnosticMap := getDiagnosticMap(readResult)
	if errors.Is(err, ErrInvalidFields) {
		return diagnosticMap, ErrInvalidFields
	}

	if err := helpers.FromMap(readResult.getConfigMap(), target); err != nil {
		return nil, fmt.Errorf("error converting map into target: %v", err)
	}

	return diagnosticMap, nil
}

func getDiagnosticMap(r readResult) map[string]string {
	if v, ok := r.(diagnosticReadResult); ok {
		return v.diagnosticMap
	}
	return nil
}
