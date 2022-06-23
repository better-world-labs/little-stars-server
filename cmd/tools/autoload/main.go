package main

import (
	"github.com/urfave/cli/v2"
	"gitlab.openviewtech.com/openview-pub/gopkg/log"
	"os"
	"path"
)

func main() {
	watchFlag := cli.BoolFlag{
		Name:  "w",
		Value: false,
		Usage: "watch files change",
	}

	app := &cli.App{
		Name:  "autoload",
		Usage: "autoload -s ${scan_package_dir} -p ${pkgName} -f ${funcName} -o ${out_dir} [-w]",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "s",
				Usage:    "scan package dir",
				Required: true,
			},

			&cli.StringFlag{
				Name:     "p",
				Value:    "",
				Usage:    "package name",
				Required: true,
			},

			&cli.StringFlag{
				Name:     "f",
				Value:    "",
				Usage:    "function name",
				Required: true,
			},

			&cli.StringFlag{
				Name:     "o",
				Value:    "",
				Usage:    "output filepath",
				Required: true,
			},

			&watchFlag,
		},
		Action: func(c *cli.Context) error {
			dirs := c.StringSlice("s")
			wd, err := os.Getwd()
			if err != nil {
				return err
			}

			for i := range dirs {
				dirs[i] = path.Join(wd, dirs[i])
			}

			loader := autoload{
				scanDir:      dirs,
				packageName:  c.String("p"),
				functionName: c.String("f"),
				outputFile:   c.String("o"),
			}
			err = loader.fillModuleInfo()
			if err != nil {
				log.Error("loader.fillModuleInfo() err:", err)
				return err
			}
			err = loader.firstGenerate()
			if err != nil {
				log.Error("loader.firstGenerate() err:", err)
				return err
			}
			if c.Bool("w") {
				log.Info("watch mode...")
				doWatch(loader.reGenerate, dirs)
			}
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err)
	}
}
