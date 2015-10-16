package plugged

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestGateway(t *testing.T) {
	examples := map[string]struct {
		name        string
		description string
		home        string
		path        string
		files       map[string]string
		scenario    [][]string
		output      string
	}{

		"default message without any plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",
			files:       map[string]string{},

			scenario: [][]string{
				{"exampleapp"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- help\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"default message with one plugin": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",

			files: map[string]string{
				"./tmp/bin/exampleapp-find": dedent(`
                                      |#!/usr/bin/env sh
                                      |arg=$1
                                      |if test "$arg" = "--plugged-description"; then
                                      |  echo -n "Find some stuff."
                                      |fi
                              `),
			},

			scenario: [][]string{
				{"exampleapp", "--plugged-install", "find"},
				{"exampleapp"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- find\t - Find some stuff.
                              |- help\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"default message with some plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",

			files: map[string]string{
				"./tmp/bin/exampleapp-find": dedent(`
                                      |#!/usr/bin/env sh
                                      |arg=$1
                                      |if test "$arg" = "--plugged-description"; then
                                      |  echo -n "Find some stuff."
                                      |fi
                              `),
				"./tmp/bin/exampleapp-activate": dedent(`
                                      |#!/usr/bin/env sh
                                      |arg=$1
                                      |if test "$arg" = "--plugged-description"; then
                                      |  echo -n "Activate stuff."
                                      |fi
                              `),
			},

			scenario: [][]string{
				{"exampleapp", "--plugged-install", "find"},
				{"exampleapp", "--plugged-install", "activate"},
				{"exampleapp"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- activate\t - Activate stuff.
                              |- find\t\t - Find some stuff.
                              |- help\t\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"default message with some plugins installed with one command": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",

			files: map[string]string{
				"./tmp/bin/exampleapp-find": dedent(`
                                      |#!/usr/bin/env sh
                                      |arg=$1
                                      |if test "$arg" = "--plugged-description"; then
                                      |  echo -n "Find some stuff."
                                      |fi
                              `),
				"./tmp/bin/exampleapp-activate": dedent(`
                                      |#!/usr/bin/env sh
                                      |arg=$1
                                      |if test "$arg" = "--plugged-description"; then
                                      |  echo -n "Activate stuff."
                                      |fi
                              `),
			},

			scenario: [][]string{
				{"exampleapp", "--plugged-install", "find", "activate"},
				{"exampleapp"},
			},

			output: dedent(`
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- activate\t - Activate stuff.
                              |- find\t\t - Find some stuff.
                              |- help\t\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"help message without any plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",
			files:       map[string]string{},

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
                              |- help\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"--help message without any plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",
			files:       map[string]string{},

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
                              |- help\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"find command when find plugin is not installed but there are some plugins": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",

			files: map[string]string{
				"./tmp/bin/exampleapp-activate": dedent(`
                                      |#!/usr/bin/env sh
                                      |echo -n "Activate stuff."
                              `),
			},

			scenario: [][]string{
				{"exampleapp", "--plugged-install", "activate"},
				{"exampleapp", "find", "stuff"},
			},

			output: dedent(`
                              |[ERROR] Unable to find plugin 'find'.
                              |Try installing it with 'exampleapp --plugged-install find'.
                              |Details: Plugin 'find' was not found
                              |
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- activate\t - Activate stuff.
                              |- help\t\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"find command when find plugin is not installed": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",
			files:       map[string]string{},

			scenario: [][]string{
				{"exampleapp", "find", "stuff"},
			},

			output: dedent(`
                              |[ERROR] Unable to find plugin 'find'.
                              |Try installing it with 'exampleapp --plugged-install find'.
                              |Details: There are no plugins installed
                              |
                              |USAGE: exampleapp command [options]
                              |
                              |exampleapp - An example CLI application.
                              |
                              |Available commands:
                              |
                              |- help\t - This info.
                              |
                              |To get help for any of commands you can do 'exampleapp help command'
                              |or 'exampleapp command --help'.
                      `),
		},

		"find command when find plugin is installed": {
			name:        "exampleapp",
			description: "An example CLI application.",
			home:        "./tmp/home",
			path:        "./tmp/bin",

			files: map[string]string{
				"./tmp/bin/exampleapp-activate": dedent(`
                                      |#!/usr/bin/env sh
                                      |echo -n "Activate stuff."
                              `),
				"./tmp/bin/exampleapp-find": dedent(`
                                      |#!/usr/bin/env sh
                                      |if test "$1" = "--plugged-description"; then
                                      |  echo -n "Activate stuff."
                                      |else
                                      |  echo "Found $1 and maybe($2)."
                                      |fi
                              `),
			},

			scenario: [][]string{
				{"exampleapp", "--plugged-install", "activate", "find"},
				{"exampleapp", "find", "stuff", "things"},
			},

			output: dedent(`
                              |Found stuff and maybe(things).
                      `),
		},
	}

	for exampleName, example := range examples {
		t.Log(exampleName)

		func() {
			if err := os.MkdirAll(example.home, 0777); err != nil {
				t.Fatalf("Unable to create home directory - %s", err)
			}
			defer os.RemoveAll(example.home)

			if err := os.MkdirAll(example.path, 0777); err != nil {
				t.Fatalf("Unable to create path directory - %s", err)
			}
			defer os.RemoveAll(example.path)

			oldPath := os.Getenv("PATH")
			os.Setenv("PATH", example.path+":"+oldPath)
			defer os.Setenv("PATH", oldPath)

			for path, contents := range example.files {
				if err := ioutil.WriteFile(path, []byte(contents), 0777); err != nil {
					t.Fatalf("Unable to create file %s - %s", path, err)
				}
			}

			stdout := &bytes.Buffer{}
			gateway := &GatewayT{
				Stdin:       bytes.NewBufferString(""),
				Stdout:      stdout,
				Home:        example.home,
				Name:        example.name,
				Description: example.description,
				ExecFn:      dumbExec(stdout),
			}

			if err := gateway.Connect(); err != nil {
				t.Fatalf("Unable to conect gateway to its store - %s", err)
			}
			defer gateway.Disconnect()
			defer os.Remove(example.home + "." + example.name + ".db")

			for _, args := range example.scenario {
				if err := gateway.Run(args); err != nil {
					t.Fatal(err)
				}
			}

			if actual := string(stdout.Bytes()); actual != example.output {
				t.Errorf(
					"\n=== Expected output ===\n%s\n=== Actual output ===\n%s\n=== END ===",
					example.output,
					actual,
				)
			}
		}()
	}
}

func dedent(s string) string {
	filter := regexp.MustCompile(`^\s*\|(.*)$`)

	s = strings.Replace(s, `\t`, "\t", -1)

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

func dumbExec(w io.Writer) func(string, []string, []string) error {
	return func(binary string, args []string, _ []string) error {
		args = args[1:]
		cmd := exec.Command(binary, args...)

		out, err := cmd.Output()
		if err != nil {
			return fmt.Errorf("Unable to execute %+v - %s", args, err)
		}

		if _, err := w.Write(out); err != nil {
			return fmt.Errorf("Unable to write output to stdout %+v - %s", args, err)
		}

		return nil
	}
}
