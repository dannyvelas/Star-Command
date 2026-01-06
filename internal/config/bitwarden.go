package config

type bitwardenConfig struct {
	BitwardenAPIURL      string `json:"bitwarden_api_url"`
	BitwardenIdentityURL string `json:"bitwarden_identity_url"`
}

var defaultBitwardenConfig = bitwardenConfig{
	BitwardenAPIURL:      "https://api.bitwarden.com",
	BitwardenIdentityURL: "https://identity.bitwarden.com",
}
