package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/dannyvelas/conflux"
	"github.com/dannyvelas/homelab/internal/app"
	"github.com/spf13/cobra"
)

func getConfigCmd() *cobra.Command {
	var targets []string

	getConfigCmd := &cobra.Command{
		Use:   "config <host-alias>",
		Short: "Generate a JSON object of configuration values for a given host",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			hostAlias := args[0]

			configs, diagnostics, err := app.GetConfig(hostAlias, targets)
			if err != nil {
				fmt.Fprintf(os.Stderr, "internal error: %s", err.Error())
				os.Exit(1)
			} else if len(diagnostics) > 0 {
				fmt.Fprintf(os.Stderr, "%v for %s:\n%s\n", conflux.ErrInvalidFields, hostAlias, conflux.DiagnosticsToTable(diagnostics))
				return
			}

			bytes, err := json.MarshalIndent(configs, "", "    ")
			if err != nil {
				fmt.Fprintf(os.Stderr, "internal error marshalling cnofigs to JSON: %s", err.Error())
				os.Exit(1)
			}

			fmt.Println(string(bytes))
		},
	}

	getConfigCmd.Flags().StringSliceVar(&targets, "for", []string{"ansible"}, "Get config for specific integration (e.g. ansible, ssh)")

	return getConfigCmd
}
