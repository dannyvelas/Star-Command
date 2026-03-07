package main

import (
	"github.com/dannyvelas/starcommand/internal/app"
	"github.com/dannyvelas/starcommand/internal/models"
	"github.com/spf13/cobra"
)

func sshAddCmd(c *models.Config, preflight *bool) *cobra.Command {
	sshAddCmd := &cobra.Command{
		Use:   "add <host>",
		Short: "Add a host to ~/.ssh/config",
		Args:  cobra.ExactArgs(1),
		RunE:  sshAddCLI(c, preflight),
	}

	return sshAddCmd
}

func sshAddCLI(c *models.Config, preflight *bool) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		hostAlias := args[0]
		return app.SSHAdd(cmd.Context(), c, hostAlias, *preflight)
	}
}
