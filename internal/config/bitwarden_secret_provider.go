package config

import (
	"fmt"

	"github.com/dannyvelas/homelab/internal/client"
)

var _ unvalidatedReader = bitwardenSecretProvider{}

type bitwardenSecretProvider struct {
	bitwardenCredProvider bitwardenCredProvider
}

func newBitwardenSecretProvider(configMap map[string]string) bitwardenSecretProvider {
	return bitwardenSecretProvider{
		bitwardenCredProvider: newBitwardenCredProvider(configMap),
	}
}

func (p bitwardenSecretProvider) ReadUnvalidated() (map[string]string, error) {
	config := newBitwardenConfig()
	if err := UnmarshalInto(p.bitwardenCredProvider, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling bitwarden creds: %v", err)
	}

	results, ok, err := validateStruct(config)
	if err != nil {
		return nil, fmt.Errorf("error validating bitwarden config: %v", err)
	} else if !ok {
		return nil, fmt.Errorf("error: invalid bitwarden configs: %s", fmtTable(results))
	}

	bitwardenClient, err := client.NewBitwardenClient(
		config.APIURL,
		config.IdentityURL,
		config.AccessToken,
		config.OrganizationID,
		config.ProjectID,
		config.StateFilePath,
	)
	if err != nil {
		return nil, fmt.Errorf("error initializing bitwarden client: %v", err)
	}

	// read bitwarden secrets
	bitwardenSecrets, err := bitwardenClient.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading bitwarden secrets: %v", err)
	}

	return bitwardenSecrets, nil
}
