package main

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ickerio/mima/parsers"
	"github.com/ickerio/mima/printer"
	"github.com/ickerio/mima/providers"
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
					conf, err := parsers.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					provider, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					server, err := provider.Info()
					if err != nil {
						return err
					}

					printer.PrintInfo(server)
					fmt.Printf("\nPassword: %v\n", server.Password)

					return nil
				},
			},
			{
				Name:  "start",
				Usage: "Starts the given server if not already online",
				Action: func(c *cli.Context) error {
					// Get config...
					conf, err := parsers.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					// Start the VPS...
					provider, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					err = provider.Start()
					if err != nil {
						return err
					}
					fmt.Println("[MIMA] Success! VPS is starting now... please wait")

					// Fetch server info...
					var server providers.Server = providers.Server{}
					fmt.Print("[MIMA] Waiting for VPS to initialize...")
					for i := 0; i < 180; i++ {
						time.Sleep(time.Second * 10)
						ser, err := provider.Info()
						if err != nil {
							return err
						}

						if ser != (providers.Server{}) && ser.Password != "not supported" && ser.Ready {
							server = ser
							break
						}

						fmt.Printf("\r[MIMA] Waiting for VPS to initialize... (%v tries)", i)
					}
					fmt.Println()
					// Ensure server is acquired.
					if (providers.Server{}) == server {
						return errors.New("Unable to fetch server info.")
					}

					service, err := parsers.GetService("minecraft.service", c.String("saves"), server.Name, server.IP, server.Password)
					if err != nil {
						return err
					}

					if err := service.Start(); err != nil {
						return err
					}
					fmt.Println("[MIMA] Server starting... please wait")

					return nil
				},
			},
			{
				Name:  "stop",
				Usage: "Stop the given server if currently online",
				Action: func(c *cli.Context) error {
					conf, err := parsers.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					provider, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					var server providers.Server = providers.Server{}
					fmt.Print("[MIMA] Waiting for VPS to initialize...")
					for i := 0; i < 180; i++ {
						ser, err := provider.Info()
						if err != nil {
							return err
						}
						if ser != (providers.Server{}) {
							server = ser
							break
						}

						fmt.Printf("\r[MIMA] Waiting for VPS to initialize... (%v tries)", i)
						time.Sleep(time.Second * 10)
					}
					fmt.Println()

					if (providers.Server{}) == server {
						return errors.New("Unable to fetch server info.")
					}

					service, err := parsers.GetService("minecraft.service", c.String("saves"), server.Name, server.IP, server.Password)
					if err != nil {
						return err
					}

					if err := service.Stop(); err != nil {
						return err
					}
					fmt.Println("[MIMA] Game server shut down.")

					err = provider.Stop()
					if err != nil {
						return err
					}
					fmt.Println("[MIMA] Success! VPS is shutting down.")

					return nil
				},
			},
			{
				Name:    "plans",
				Aliases: []string{"plan", "p"},
				Usage:   "Lists all the plans of a particular service",
				Action: func(c *cli.Context) error {
					provider, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					plans, err := provider.Plans()
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
					provider, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					regions, err := provider.Regions()
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
					provider, err := providers.GetNoAuth(c.Args().Get(0))
					if err != nil {
						return err
					}

					os, err := provider.OS()
					if err != nil {
						return err
					}
					printer.PrintOS(os)

					return nil
				},
			},
			{
				Name:  "create",
				Usage: "Start the VPS, creating a new server.",
				Action: func(c *cli.Context) error {
					conf, err := parsers.GetConfig(c.String("config"))
					if err != nil {
						return err
					}

					// Start the VPS...
					provider, err := providers.GetFromConfig(conf, c.Args().Get(0))
					if err != nil {
						return err
					}

					err = provider.Start()
					if err != nil {
						return err
					}
					fmt.Println("[MIMA] Success! VPS is starting now... please wait")

					// Fetch server info...
					var server providers.Server = providers.Server{}
					fmt.Print("[MIMA] Fetching server info...")
					for i := 0; i < 36; i++ {
						time.Sleep(time.Second * 10)
						ser, err := provider.Info()
						if err != nil {
							return err
						}

						if ser != (providers.Server{}) && ser.Password != "not supported" && ser.Ready {
							server = ser
							break
						}

						fmt.Printf("\r[MIMA] Fetching server info - %v tries...", i)
					}
					fmt.Println()
					// Ensure server is acquired.
					if (providers.Server{}) == server {
						return errors.New("Unable to fetch server info.")
					}

					service, err := parsers.GetService("minecraft.service", c.String("saves"), server.Name, server.IP, server.Password)
					if err != nil {
						return err
					}

					if err := service.Create(); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("[MIMA ERROR] " + err.Error())
	}
}
