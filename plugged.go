// Package plugged is a library for writing extendable CLI applications.
package plugged

import (
	"fmt"
	"io"

	"github.com/boltdb/bolt"
)

var builtinHandlers = map[string]actionHandler{
	"help":              actionHandler((*GatewayT).helpAction),
	"--plugged-install": actionHandler((*GatewayT).installAction),
}

type actionHandler func(g *GatewayT, action string, args []string) error

// GatewayT represents a "Gateway" CLI application configuration.
type GatewayT struct {
	Stdin       io.Reader
	Stdout      io.Writer
	Home        string
	Name        string
	Description string

	store *bolt.DB
}

// Run is for executing a command according to provided arguments.
func (g *GatewayT) Run(args []string) error {
	action, args := argsToAction(args)

	if handler, ok := builtinHandlers[action]; ok {
		if err := handler(g, action, args); err != nil {
			return err
		}
	}

	return nil
}

func (g *GatewayT) Connect() error {
	var err error

	g.store, err = bolt.Open(g.Home+"/."+g.Name+".db", 0600, nil)
	if err != nil {
		return fmt.Errorf("Unable to connect to embedded database - %s", err)
	}

	return nil
}

func (g *GatewayT) Disconnect() {
	g.store.Close()
}

func (g *GatewayT) Plugins() ([]*pluginT, error) {
	var plugins []*pluginT

	err := g.store.View(func(tx *bolt.Tx) error {
		var err error

		b := tx.Bucket([]byte("plugins"))
		if b == nil {
			return nil
		}

		if plugins, err = listPlugins(b); err != nil {
			return fmt.Errorf("Unable to get plugins - %s", err)
		}

		return nil
	})

	return plugins, err
}

func (g *GatewayT) updatePlugin(p *pluginT) error {
	return g.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("plugins"))
		if err != nil {
			return fmt.Errorf("Unable to obtain bucket 'plugins' - %s", err)
		}

		if err := p.save(b); err != nil {
			return fmt.Errorf("Unable to save plugin to bucket 'plugins' - %s", err)
		}

		return nil
	})
}

func (g *GatewayT) helpAction(string, []string) error {
	plugins, err := g.Plugins()
	if err != nil {
		return err
	}

	commandList := &commandListView{
		PluginList: plugins,
	}

	availableCommands, err := commandList.render()
	if err != nil {
		return err
	}

	help := &helpView{
		Name:              g.Name,
		Description:       g.Description,
		AvailableCommands: availableCommands,
	}

	if err := help.render(g.Stdout); err != nil {
		return err
	}
	return nil
}

func (g *GatewayT) installAction(_ string, plugins []string) error {
	for _, name := range plugins {
		p := newPlugin(g.Name, name)

		if err := p.install(g); err != nil {
			fmt.Printf("%s: Failed to get metadata - %s\n", name, err)
		}
	}

	return nil
}

func argsToAction(args []string) (string, []string) {
	if len(args) == 1 || args[1] == "--help" {
		return "help", []string{}
	}

	return args[1], args[2:]
}
