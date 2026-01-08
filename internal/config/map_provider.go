package config

import (
	"encoding/json"
	"fmt"
)

var _ provider = mapProvider{}

type mapProvider struct {
	mapAsBytes []byte
}

func newMapProvider(m map[string]string) (mapProvider, error) {
	bytes, err := json.Marshal(m)
	if err != nil {
		return mapProvider{}, fmt.Errorf("error marshalling map: %v", err)
	}

	return mapProvider{mapAsBytes: bytes}, nil
}

func (p mapProvider) UnmarshalInto(target any) error {
	if err := json.Unmarshal(p.mapAsBytes, target); err != nil {
		return fmt.Errorf("error unmarshalling map into target: %v", err)
	}
	return nil
}
