package plugged

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/boltdb/bolt"
)

type pluginT struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
	AppName     string `json:"AppName"`
}

func newPlugin(appName, name string) *pluginT {
	return &pluginT{
		Name:    name,
		AppName: appName,
	}
}

func listPlugins(store *bolt.Bucket) ([]*pluginT, error) {
	plugins := []*pluginT{}

	err := store.ForEach(func(key, data []byte) error {
		plugin, err := decodePlugin(data, string(key))
		if err != nil {
			fmt.Printf("[ERROR] %s\n", err)
		}

		plugins = append(plugins, plugin)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch plugins from store - %s", err)
	}

	return plugins, nil
}

func pluginFrom(store *bolt.Bucket, name string) (*pluginT, error) {
	data := store.Get([]byte(name))
	if data == nil {
		return nil, fmt.Errorf("Plugin '%s' was not found", name)
	}

	return decodePlugin(data, name)
}

func decodePlugin(data []byte, name string) (*pluginT, error) {
	plugin := &pluginT{}

	if err := json.Unmarshal(data, plugin); err != nil {
		return nil, fmt.Errorf("Unable to unmarshal plugin %s data %+v - %s", name, data, err)
	}

	return plugin, nil
}

func (p *pluginT) install(g *GatewayT) error {
	cmdName := p.AppName + "-" + p.Name
	cmd := exec.Command(cmdName, "--plugged-description")

	description, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("'%s --plugged-description' returned an error - %s", cmdName, err)
	}

	p.Description = string(description)

	if err := g.updatePlugin(p); err != nil {
		return fmt.Errorf("Unable to save plugin to storage - %s", err)
	}
	return nil
}

func (p *pluginT) save(store *bolt.Bucket) error {
	data, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Unable to marshal plugin %+v to json - %s", *p, err)
	}

	return store.Put([]byte(p.Name), data)
}

func (p *pluginT) run(execFn func(string, []string, []string) error, args []string) error {
	cmdName := p.AppName + "-" + p.Name

	binary, err := exec.LookPath(cmdName)
	if err != nil {
		return fmt.Errorf("Unable to find binary for plugin '%s' - %s", p.Name, err)
	}

	args = append([]string{cmdName}, args...)

	if err := execFn(binary, args, os.Environ()); err != nil {
		return fmt.Errorf("Plugin '%s' failed to exec with args: %+v - %s", p.Name, args, err)
	}

	return nil
}
