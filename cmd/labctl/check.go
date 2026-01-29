package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/dannyvelas/conflux"
	"github.com/dannyvelas/homelab/internal/app"
	"github.com/dannyvelas/homelab/internal/helpers"
	"github.com/spf13/cobra"
)

func checkCmd() *cobra.Command {
	checkCmd := &cobra.Command{
		Use:   "check <host-alias> target1 [targets]",
		Short: "Print a diagnostic report of all the configs that were found/missing for a given resource",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			hostAlias := args[0]
			configMux := conflux.NewConfigMux(
				conflux.WithYAMLFileReader(helpers.FallbackFile, conflux.WithPath(helpers.GetConfigPath(hostAlias))),
				conflux.WithEnvReader(),
				conflux.WithBitwardenSecretReader(),
			)

			targets, err := toTargets(args[1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			diagnostics, err := app.Check(configMux, hostAlias, targets)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s\n", err.Error())
				os.Exit(1)
			}

			fmt.Printf("Configs needed for host(%s):\n%s\n", hostAlias, app.DiagnosticsToTable(diagnostics))
		},
	}

	return checkCmd
}

var m = map[string]token{
	"terraform": terraform,
	"ssh":       ssh,
	"check":     check,

	"ansible":   ansible,
	"inventory": inventory,
	"playbook":  playbook,

	"add":   add,
	"run":   run,
	"apply": apply,
}

type token int

const (
	nilToken token = iota

	terraform
	ssh
	check

	ansible
	inventory
	playbook

	add
	run
	apply

	eof
)

func toTargets(args []string) ([]app.Target, error) {
	for _, arg := range args {
		tokens := scan(arg)
		target, err := parse(tokens)
		fmt.Println(target, err)
	}
	return nil, errors.New("unimplemented")
}

func scan(source string) chan token {
	tokens := make(chan token)
	go func() {
		start, current := 0, 0
		for current < len(source) {
			start = current
			newCurrent, token, err := scanToken(source, start, current)
			current = newCurrent
			if errors.Is(err, errSkip) {
				continue
			} else if err != nil {
				fmt.Printf("error: %v\n", err)
				continue
			}
			tokens <- token
		}

		tokens <- eof
	}()
	return tokens
}

var errSkip = fmt.Errorf("skip token")

func scanToken(source string, start, current int) (int, token, error) {
	newCurrent, c := current+1, source[current]

	if c == ':' {
		return newCurrent, nilToken, errSkip
	}

	if !isLower(c) {
		return newCurrent, nilToken, fmt.Errorf("invalid token")
	}

	for newCurrent < len(source) && isLower(source[newCurrent]) {
		newCurrent += 1
	}

	lexeme := source[start:newCurrent]
	tok, ok := m[lexeme]
	if !ok {
		return newCurrent, nilToken, fmt.Errorf("invalid token")
	}

	return newCurrent, tok, nil
}

func isLower(b byte) bool {
	return b >= 'a' && b <= 'z'
}

func parse(tokens chan token) (app.Target, error) {
	resource, err := parseResource(tokens)
	if err != nil {
		return app.Target{}, fmt.Errorf("error parsing resource: %v", err)
	}
	action, err := parseAction(tokens)
	if err != nil {
		return app.Target{}, fmt.Errorf("error parsing action: %v", err)
	}

	return app.Target{Resource: resource, Action: action}, nil
}

func parseResource(tokens chan token) (app.Resource, error) {
	switch <-tokens {
	case ansible:
		switch <-tokens {
		case playbook:
			return app.AnsiblePlaybookResource, nil
		case inventory:
			return app.AnsibleInventoryResource, nil
		default:
			return "", fmt.Errorf("invalid resource")
		}
	case ssh:
		return app.SSHResource, nil
	case terraform:
		return app.TerraformResource, nil
	default:
		return "", fmt.Errorf("invalid resource")
	}
}

func parseAction(tokens chan token) (app.Action, error) {
	switch <-tokens {
	case run:
		return app.RunAction, nil
	case add:
		return app.AddAction, nil
	case apply:
		return app.ApplyAction, nil
	default:
		return "", fmt.Errorf("invalid resource")
	}
}
