package env

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	BitwardenAccessToken    string
	BitwardenOrganizationID string
	BitwardenProjectID      string
	BitwardenStateFilePath  string
}

// New returns an Env struct of all expected environmental variables
func New() Env {
	godotenv.Load()

	return Env{
		BitwardenAccessToken:    os.Getenv("BWS_ACCESS_TOKEN"),
		BitwardenOrganizationID: os.Getenv("BWS_PROJECT_ID"),
		BitwardenProjectID:      os.Getenv("BWS_ORGANIZATION_ID"),
		BitwardenStateFilePath:  os.Getenv("BWS_STATE_FILE_PATH"),
	}
}
