package plugged

import (
	"os"
	"syscall"
)

// Gateway creates a main "Gateway" style application
func Gateway(name, description string, args []string) {
	gateway := &GatewayT{
		Stdin:       os.Stdin,
		Stdout:      os.Stdout,
		Home:        os.Getenv("HOME"),
		Name:        name,
		Description: description,
		ExecFn:      syscall.Exec,
	}

	gateway.Run(args)
}
