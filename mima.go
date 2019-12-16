package main

import (
	"fmt"
	"os"

	"github.com/ickerio/mima/providers"
	"github.com/ickerio/mima/util"
	"github.com/urfave/cli/v2"
)

func main() {
	var (
		conf util.Config
		prov providers.Provider
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
			configuration, err := util.GetConfig(c.String("config"))
			if err != nil {
				return err
			}
			conf = configuration

			if c.Args().Present() {
				provider, err := providers.Get(conf, c.Args().Get(1))
				if err != nil {
					return err
				}
				prov = provider
			}

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Displays info on the server",
				Action: func(c *cli.Context) error {
					ser, err := prov.Info()
					if err != nil {
						return err
					}

					fmt.Printf(
						"%v running %v in %v at %v\n%v memory, %v storage, %v cpus",
						ser.Name, ser.Os, ser.Location, ser.IP, ser.Memory, ser.Storage, ser.CPUCount,
					)

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
