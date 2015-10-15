// Package plugged is a library for writing extendable CLI applications.
package plugged

import (
	"fmt"
	"io"
	"os"
	"text/template"
)

var (
	gatewayHelpTemplate = template.Must(template.New("gatewayHelpView").Parse(
		`USAGE: {{.Name}} command [options]

{{.Name}} - {{.Description}}

Available commands:

- help - This info.

To get help for any of commands you can do '{{.Name}} help command'
or '{{.Name}} command --help'.
`,
	))
)

// GatewayT represents a "Gateway" CLI application configuration.
type GatewayT struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Name        string
	Description string
}

// Run is for executing a command according to provided arguments.
func (g *GatewayT) Run(args []string) {
	if err := g.renderHelpView(); err != nil {
		fmt.Printf("Unexpected error occurred - %s", err)
		os.Exit(1)
	}
}

func (g *GatewayT) renderHelpView() error {
	if err := gatewayHelpTemplate.Execute(g.Stdout, g); err != nil {
		return fmt.Errorf("Unable to execute helpView template on %v - %s", g, err)
	}
	return nil
}
