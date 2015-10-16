// Package plugged is a library for writing extendable CLI applications.
package plugged

import (
	"fmt"
	"io"

	"github.com/boltdb/bolt"
)

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
	if len(args) == 1 || args[1] == "help" || args[1] == "--help" {
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

	if args[1] == "--plugged-install" {
		plugins := args[2:]
		for _, name := range plugins {
			p := NewPlugin(g.Name, name)

			if err := p.Install(g); err != nil {
				fmt.Printf("%s: Failed to get metadata - %s\n", name, err)
			}
		}

		return nil
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

func (g *GatewayT) UpdatePlugin(p *PluginT) error {
	return g.store.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("plugins"))
		if err != nil {
			return fmt.Errorf("Unable to obtain bucket 'plugins' - %s", err)
		}

		if err := p.Save(b); err != nil {
			return fmt.Errorf("Unable to save plugin to bucket 'plugins' - %s", err)
		}

		return nil
	})
}

func (g *GatewayT) Plugins() ([]*PluginT, error) {
	var plugins []*PluginT

	err := g.store.View(func(tx *bolt.Tx) error {
		var err error

		b := tx.Bucket([]byte("plugins"))
		if b == nil {
			return nil
		}

		if plugins, err = ListPlugins(b); err != nil {
			return fmt.Errorf("Unable to get plugins - %s", err)
		}

		return nil
	})

	return plugins, err
}
