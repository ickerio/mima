package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ickerio/mima/config"
	"github.com/ickerio/mima/providors"
	"github.com/urfave/cli"
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

			fmt.Println(c.Args().First())
			providor, err := providors.Get(conf, c.Args().First())
			if err != nil {
				return err
			}
			prov = providor

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Displays info on given game server",
				Action: func(c *cli.Context) error {
					fmt.Println(conf)
					fmt.Println(prov)
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
