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
			Name:  "b",
			Value: ".",
			Usage: "-b ${baseDir}",
		},
		&cli.StringFlag{
			Name:  "e",
			Value: "local",
			Usage: "-e ${environment}",
		},
	}

	app := &cli.App{
		Flags:  flags,
		Action: command.Run,
	}

	app.Name = "Openviewtech AED Api Service"
	app.Version = "1.10.0"
	app.Copyright = "(c) 2022 openviewtech"
	app.EnableBashCompletion = true

	if err := app.Run(os.Args); err != nil {
		fmt.Println(err)
	}
}
