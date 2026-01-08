package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

var _ unvalidatedReader = fileProvider{}

type fileProvider struct {
	hostName string
	verbose  bool
}

func newFileProvider(hostName string, verbose bool) fileProvider {
	return fileProvider{
		hostName: hostName,
		verbose:  verbose,
	}
}

func (p fileProvider) ReadUnvalidated() (map[string]string, error) {
	m := make(map[string]string)
	hostConfigFile := filepath.Join(configDir, fmt.Sprintf("%s.yml", p.hostName))
	for _, file := range []string{fallbackConfigFile, hostConfigFile} {
		data, err := os.ReadFile(file)
		if errors.Is(err, os.ErrNotExist) {
			if p.verbose {
				fmt.Fprintf(os.Stderr, "warning: %s config file not found\n", file)
			}
			continue
		} else if err != nil {
			return nil, fmt.Errorf("error reading config file(%s): %v", file, err)
		}
		if err := yaml.Unmarshal(data, m); err != nil {
			return nil, fmt.Errorf("error unmarshalling config file (%s): %v", file, err)
		}
	}
	return m, nil
}
