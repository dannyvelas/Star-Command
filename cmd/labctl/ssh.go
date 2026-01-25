package main

import "github.com/spf13/cobra"

func sshCmd() *cobra.Command {
	sshCmd := &cobra.Command{
		Use:   "ssh",
		Short: "Create a resource",
	}

	sshCmd.AddCommand(sshAddCmd())

	return sshCmd
}
