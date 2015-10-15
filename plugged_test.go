package plugged

import (
	"bytes"
	"regexp"
	"strings"
	"testing"
)

func TestGateway(t *testing.T) {
	examples := map[string]struct {
		name        string
		description string
		scenario    [][]string
		output      string
	}{

		"help message without any plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			scenario: [][]string{
				{"exampleapp", "help"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- help - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"--help message without any plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			scenario: [][]string{
				{"exampleapp", "--help"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- help - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},
	}

	for exampleName, example := range examples {
		t.Log(exampleName)

		stdout := &bytes.Buffer{}
		gateway := &GatewayT{
			Stdin:       bytes.NewBufferString(""),
			Stdout:      stdout,
			Name:        example.name,
			Description: example.description,
		}

		for _, args := range example.scenario {
			gateway.Run(args)
		}

		if actual := string(stdout.Bytes()); actual != example.output {
			t.Errorf(
				"\n=== Expected output ===\n%s\n=== Actual output ===\n%s\n=== END ===",
				example.output,
				actual,
			)
		}
	}
}

func dedent(s string) string {
	filter := regexp.MustCompile(`^\s*\|(.*)$`)

	acc := ""
	for _, line := range strings.Split(s, "\n") {
		matches := filter.FindStringSubmatch(line)
		if len(matches) < 2 {
			continue
		}

		dedented := matches[1]
		acc += dedented + "\n"
	}
	return acc
}
