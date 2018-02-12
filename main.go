package main

import (
	"fmt"
	"github.com/dcarbone/cs-zone-cloner/command"
	"github.com/dcarbone/cs-zone-cloner/command/backup"
	"github.com/dcarbone/cs-zone-cloner/command/restore"
	"github.com/mitchellh/cli"
	stdlog "log"
	"os"
)

func main() {

	l := command.NewMutableLogger(stdlog.New(os.Stderr, "", stdlog.LstdFlags))

	c := cli.NewCLI("cs-zone-cloner", "dev")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"backup": func() (cli.Command, error) {
			return backup.New(os.Args[0], l), nil
		},
		"restore": func() (cli.Command, error) {
			return restore.New(os.Args[0], l), nil
		},
	}

	status, err := c.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s\n", err)
	}

	os.Exit(status)
}
