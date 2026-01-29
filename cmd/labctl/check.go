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
		tokens, errors := scan(arg)
		fmt.Println(tokens)
		fmt.Println(errors)
		target, err := parseTarget(tokens)
		fmt.Println(target, err)
	}
	return nil, errors.New("unimplemented")
}

func scan(source string) ([]token, []error) {
	errs := make([]error, 0)
	tokens := make([]token, 0)
	start, current := 0, 0
	for current < len(source) {
		start = current
		newCurrent, token, err := scanToken(source, start, current)
		current = newCurrent
		if errors.Is(err, errSkip) {
			continue
		} else if err != nil {
			errs = append(errs, err)
			continue
		}
		tokens = append(tokens, token)
	}

	tokens = append(tokens, eof)
	return tokens, errs
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

func parseTarget(tokens []token) (app.Target, error) {
	newCurrent, resource, err := parseResource(0, tokens)
	if err != nil {
		return app.Target{}, fmt.Errorf("error parsing resource: %v", err)
	}
	action, err := parseAction(newCurrent, tokens)
	if err != nil {
		return app.Target{}, fmt.Errorf("error parsing action: %v", err)
	}

	return app.Target{Resource: resource, Action: action}, nil
}

func parseResource(current int, tokens []token) (int, app.Resource, error) {
	if tokens[current] == ansible {
		newCurrent := advance(current, tokens)
		if tokens[newCurrent] == playbook {
			return advance(newCurrent, tokens), app.AnsiblePlaybookResource, nil
		} else if tokens[newCurrent] == inventory {
			return advance(newCurrent, tokens), app.AnsibleInventoryResource, nil
		} else {
			return advance(newCurrent, tokens), "", fmt.Errorf("invalid resource")
		}
	} else if tokens[current] == ssh {
		return advance(current, tokens), app.SSHResource, nil
	} else if tokens[current] == terraform {
		return advance(current, tokens), app.TerraformResource, nil
	} else {
		return current, "", fmt.Errorf("invalid resource")
	}
}

func parseAction(current int, tokens []token) (app.Action, error) {
	switch tokens[current] {
	if tokens[current] == run {
		return app.RunAction, nil
	} else if tokens[current] == add {
		return app.AddAction, nil
	} else if tokens[current] == apply {
		return app.ApplyAction, nil
	} else {
		return "", fmt.Errorf("invalid resource")
	}
	}
}

//	func match(current int, tokens []token, token token) (int, bool) {
//		if tokens[current] == token {
//			newCurrent := advance(current, tokens)
//			return newCurrent, true
//		}
//		return current, false
//	}
func advance(current int, tokens []token) int {
	if tokens[current] == eof {
		return current
	}
	return current + 1
}
