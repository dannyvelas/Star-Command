package main

import (
	"github.com/dannyvelas/starcommand/internal/models"
	"github.com/spf13/cobra"
)

func terraformCmd(c *models.Config, preflight *bool) *cobra.Command {
	terraformCmd := &cobra.Command{
		Use:   "terraform",
		Short: "Execute terraform commands",
	}

	terraformCmd.AddCommand(terraformApplyCmd(c, preflight))

	return terraformCmd
}
