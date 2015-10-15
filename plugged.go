package plugged

import (
	"fmt"
	"io"
	"text/template"
)

type GatewayT struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Name        string
	Description string
}

func (g *GatewayT) Run(args []string) {
	g.renderHelpView()
}

func (g *GatewayT) renderHelpView() {
	tmpl, err := template.New("gatewayHelpView").Parse(
		`USAGE: {{.Name}} command [options]

{{.Name}} - {{.Description}}

Available commands:

- help - This info.

To get help for any of commands you can do '{{.Name}} help command'
or '{{.Name}} command --help'.
`,
	)
	if err != nil {
		panic(fmt.Errorf("Unable to parse helpView template: %s", err))
	}

	if err := tmpl.Execute(g.Stdout, g); err != nil {
		panic(fmt.Errorf("Unable to execute helpView template on %v: %s", g, err))
	}
}
