package main

import (
	"fmt"
	"os"

	"github.com/ickerio/mima/printer"
	"github.com/ickerio/mima/providers"
	"github.com/ickerio/mima/util"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "mima",
		Usage: "Server manager",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "config, c",
				Usage: "Load configuration from `FILE`",
				Value: ".mima.yml",
			},
			&cli.StringFlag{
				Name:  "saves, s",
				Usage: "`DIR`ectory of saved data",
				Value: "saves",
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Displays info on the server",
				Action: func(c *cli.Context) error {
					conf, err := util.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					prov, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					ser, err := prov.Info()
					if err != nil {
						return err
					}

					printer.PrintInfo(ser)

					return nil
				},
			},
			{
				Name:  "start",
				Usage: "Starts the given server if not already online",
				Action: func(c *cli.Context) error {
					conf, err := util.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					prov, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					err = prov.Start()
					if err != nil {
						return err
					}

					fmt.Println("Success! VPS is starting now... please wait")

					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "Stop the given server if currently online",
				Action: func(c *cli.Context) error {
					conf, err := util.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					prov, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					err = prov.Stop()
					if err != nil {
						return err
					}

					fmt.Println("Success! VPS is shutting down")
					return nil
				},
			},
			{
				Name:    "plans",
				Aliases: []string{"plan", "p"},
				Usage:   "Lists all the plans of a particular service",
				Action: func(c *cli.Context) error {
					prov, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					plans, err := prov.Plans()
					if err != nil {
						return err
					}
					printer.PrintPlans(plans)

					return nil
				},
			},
			{
				Name:    "regions",
				Aliases: []string{"region", "r"},
				Usage:   "Lists all the regions of a particular service",
				Action: func(c *cli.Context) error {
					prov, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					regions, err := prov.Regions()
					if err != nil {
						return err
					}
					printer.PrintRegions(regions)

					return nil
				},
			},
			{
				Name:  "os",
				Usage: "Lists all the operating systems of a particular service",
				Action: func(c *cli.Context) error {
					prov, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					os, err := prov.OS()
					if err != nil {
						return err
					}
					printer.PrintOS(os)

					return nil
				},
			},
			{
				Name:    "test",
				Aliases: []string{"test", "t"},
				Usage:   "Test code",
				Action: func(c *cli.Context) error {

					//services.CopyFile("README.md", "./NEW.md", "45.77.234.221", "root", "n7]CcTqhv-sZ_u[9")
					//services.SaveFile("/root/README.md", "saved.md", "45.77.234.221", "root", "n7]CcTqhv-sZ_u[9")

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
