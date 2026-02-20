package main

import (
	"fmt"

	"github.com/dannyvelas/conflux"
	"github.com/dannyvelas/homelab/internal/app"
	"github.com/spf13/cobra"
)

func ansiblePlaybookCmd(configMux *conflux.ConfigMux, preflight bool) []*cobra.Command {
	playbooks := []string{"bootstrap-server", "setup-host", "setup-vm"}
	commands := make([]*cobra.Command, 0, len(playbooks))

	for _, playbook := range playbooks {
		command := &cobra.Command{
			Use:   playbook,
			Short: fmt.Sprintf("Run the %s ansible playbook", playbook),
			Args:  cobra.ExactArgs(1),
			RunE:  run(configMux, preflight, playbook),
		}

		commands = append(commands, command)
	}
	return commands
}

func run(configMux *conflux.ConfigMux, preflight bool, playbook string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		diagnostics, err := app.AnsibleRun(ctx, configMux, preflight, playbook)
		if err != nil {
			return err
		}

		if len(diagnostics) > 0 {
			fmt.Printf("%s\n", app.DiagnosticsToTable(diagnostics))
		}
		return nil
	}
}
