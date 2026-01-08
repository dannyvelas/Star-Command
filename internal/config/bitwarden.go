package config

type bitwardenConfig struct {
	APIURL         string `json:"bitwarden_api_url"`
	IdentityURL    string `json:"bitwarden_identity_url"`
	AccessToken    string `json:"bitwarden_access_token"`
	OrganizationID string `json:"bitwarden_organization_id"`
	ProjectID      string `json:"bitwarden_project_id"`
	StateFilePath  string `json:"bitwarden_state_file_path"`
}

func newBitwardenConfig() bitwardenConfig {
	return bitwardenConfig{
		APIURL:        "https://api.bitwarden.com",
		IdentityURL:   "https://identity.bitwarden.com",
		StateFilePath: ".bw_state",
	}
}
