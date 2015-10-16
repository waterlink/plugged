package plugged

import (
	"encoding/json"
	"fmt"
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

	err := store.ForEach(func(_, data []byte) error {
		plugin := &pluginT{}

		if err := json.Unmarshal(data, plugin); err != nil {
			fmt.Printf("Unable to unmarshal plugin data '%s' - %s, ignoring", data, err)
		}

		plugins = append(plugins, plugin)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Unable to fetch plugins from store - %s", err)
	}

	return plugins, nil
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
