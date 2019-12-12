package main

import (
	"fmt"
	"os"

	"github.com/ickerio/mima/config"
	"github.com/ickerio/mima/providors"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		conf config.Config
		prov providors.Providor
	)

	app := &cli.App{
		Name:  "mima",
		Usage: "Server manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Usage: "Load configuration from `FILE`",
				Value: ".mima.yml",
			},
		},
		Before: func(c *cli.Context) error {
			configuration, err := config.Get(c.String("config"))
			if err != nil {
				return err
			}
			conf = configuration

			if c.Args().Present() {
				providor, err := providors.Get(conf, c.Args().Get(1))
				if err != nil {
					return err
				}
				prov = providor
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Displays info on given game server",
				Action: func(c *cli.Context) error {
					res, err := prov.Info()
					if err != nil {
						fmt.Println(err)
					} else {
						fmt.Println(res)
					}
					return nil
				},
			},
			{
				Name:    "start",
				Aliases: []string{"s"},
				Usage:   "Starts the given server if not already online",
				Action: func(c *cli.Context) error {
					fmt.Printf("start %q", c.Args().Get(0))
					return nil
				},
			},
			{
				Name:    "end",
				Aliases: []string{"e"},
				Usage:   "Stop the given server if currently online",
				Action: func(c *cli.Context) error {
					fmt.Printf("end %q", c.Args().Get(0))
					return nil
				},
			},
		},
	}

	app.Run(os.Args)
}
