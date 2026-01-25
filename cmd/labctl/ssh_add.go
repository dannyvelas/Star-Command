package main

import (
	"fmt"
	"os"

	"github.com/dannyvelas/conflux"
	"github.com/dannyvelas/homelab/internal/handlers"
	"github.com/dannyvelas/homelab/internal/helpers"
	"github.com/spf13/cobra"
)

func sshAddCmd() *cobra.Command {
	var targets []string

	sshAddCmd := &cobra.Command{
		Use:       "add <host-alias>",
		ValidArgs: handlers.GetSupportedHostAliases(),
		Short:     "Update the `~/.ssh/config` file to connect to a given host",
		Args:      cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hostAlias := args[0]
			configMux := conflux.NewConfigMux(
				conflux.WithYAMLFileReader(helpers.FallbackFile, conflux.WithPath(helpers.GetConfigPath(hostAlias))),
				conflux.WithEnvReader(),
				conflux.WithBitwardenSecretReader(),
			)

			handler, err := handlers.New(configMux, hostAlias, targets)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			diagnostics, err := handler.SetFile()
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			for _, diagnostic := range diagnostics {
				fmt.Printf("- %s\n", diagnostic)
			}

			fmt.Println("SSH config updated successfully!")
		},
	}

	sshAddCmd.Flags().StringSliceVar(&targets, "for", []string{"ssh"}, "Write or append to the corresponding file")

	return sshAddCmd
}
