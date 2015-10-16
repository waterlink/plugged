package plugged

import (
	"bytes"
	"fmt"
	"io"
	"text/tabwriter"
	"text/template"
)

var gatewayHelpTemplate = template.Must(template.New("gatewayHelpView").Parse(
	`USAGE: {{.Name}} command [options]

{{.Name}} - {{.Description}}

Available commands:
{{.AvailableCommands}}

To get help for any of commands you can do '{{.Name}} help command'
or '{{.Name}} command --help'.
`,
))

type helpView struct {
	Name              string
	Description       string
	AvailableCommands string
}

func (v *helpView) render(w io.Writer) error {
	if err := gatewayHelpTemplate.Execute(w, v); err != nil {
		return fmt.Errorf("Unable to execute helpView template on %v - %s", v, err)
	}
	return nil
}

var commandListTemplate = template.Must(template.New("commandListView").Parse(
	"{{range .PluginList}}\n- {{.Name}}\t - {{.Description}}{{end}}\n- help\t - This info.",
))

type commandListView struct {
	PluginList []*pluginT
}

func (v *commandListView) render() (string, error) {
	buf := &bytes.Buffer{}
	w := tabwriter.NewWriter(buf, 0, 4, 0, '\t', 0)

	if err := commandListTemplate.Execute(w, v); err != nil {
		return "", fmt.Errorf("Unable to execute commandList template on %v - %s", v, err)
	}

	if err := w.Flush(); err != nil {
		return "", fmt.Errorf("Unable to flush tabwriter - %s", err)
	}

	return string(buf.Bytes()), nil
}
