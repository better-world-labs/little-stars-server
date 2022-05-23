package main

import (
	"aed-api-server/internal/command"
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "c",
			Value: "config-local.yaml",
			Usage: "config file name",
		},
	}

	app := &cli.App{
		Flags:  flags,
		Action: command.Run,
	}

	app.Name = "Openviewtech AED Api Service"
	app.Usage = "DID & VC Repository"
	app.Version = "1.2.1"
	app.Copyright = "(c) 2021 openviewtech"
	app.EnableBashCompletion = true
	app.Commands = []*cli.Command{
		{
			Name:   "run",
			Flags:  flags,
			Action: command.Run,
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
