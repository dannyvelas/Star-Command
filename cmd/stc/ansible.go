package main

import (
	"github.com/dannyvelas/starcommand/internal/models"
	"github.com/spf13/cobra"
)

func ansibleCmd(c *models.Config, preflight *bool) *cobra.Command {
	ansibleCmd := &cobra.Command{
		Use:   "ansible",
		Short: "Execute ansible commands",
	}

	ansibleCmd.AddCommand(ansiblePlaybookCmd(c, preflight)...)

	return ansibleCmd
}
