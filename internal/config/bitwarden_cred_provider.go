package config

var _ unvalidatedReader = bitwardenCredProvider{}

type bitwardenCredProvider struct {
	configMap map[string]string
}

func newBitwardenCredProvider(configMap map[string]string) bitwardenCredProvider {
	return bitwardenCredProvider{
		configMap: configMap,
	}
}

func (p bitwardenCredProvider) ReadUnvalidated() (map[string]string, error) {
	return p.configMap, nil
}
