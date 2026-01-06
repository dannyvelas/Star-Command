package main

import (
	"github.com/dannyvelas/homelab/internal/env"
)

func main() {
	envVars := env.New()

	initialize(envVars)
	execute()
}
