package app

import (
	"errors"
	"strings"
	"testing"
)

func TestTransformTerraformVersion_Success(t *testing.T) {
	desiredVersion := "~> 1.13.3"

	// Input HCL that already has that exact version
	inputHCL := []byte(`terraform {
  required_version = "~> 1.13.3"
}

provider "proxmox" {
  endpoint = "https://10.0.0.1:8006/api2/json"
}`)

	output, err := transformTerraformVersion(inputHCL, "test.tf", desiredVersion)

	if !errors.Is(err, errAlreadyExists) {
		t.Fatalf("expected error %v, but got %v", errAlreadyExists, err)
	}

	if output != nil {
		t.Errorf("expected nil output when version already exists, but got %d bytes", len(output))
	}
}

func TestTransformTerraformVersion_Error(t *testing.T) {
	v := "~> 1.13.3"

	tests := []struct {
		name          string
		inputHCL      string
		containsInOut string // Snippet we expect in the output
	}{
		{
			name:     "No terraform block",
			inputHCL: `provider "proxmox" {}`,
			containsInOut: `terraform {
  required_version = "~> 1.13.3"
}`,
		},
		{
			name: "Terraform block exists but no version field",
			inputHCL: `terraform {
  required_providers {}
}`,
			containsInOut: `required_version = "~> 1.13.3"`,
		},
		{
			name: "Version exists but is different",
			inputHCL: `terraform {
  required_version = "1.12.0"
}`,
			containsInOut: `required_version = "~> 1.13.3"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := transformTerraformVersion([]byte(tt.inputHCL), "test.tf", v)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Check if output contains what we expected
			if !strings.Contains(string(output), tt.containsInOut) {
				t.Errorf("output missing expected content.\nGot:\n%s\nExpected to contain: %s", string(output), tt.containsInOut)
			}
		})
	}
}
