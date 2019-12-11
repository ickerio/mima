package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

type Config struct {
	VultrKey string `yaml:"vultr_key"`
	DOKey    string `yaml:"do_key"`
	Servers  []struct {
		Name     string `yaml:"name"`
		Providor string `yaml:"providor"`
		Plan     string `yaml:"plan"`
	} `yaml:"servers"`
}

func getConfig(fileName string) Config {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	decoder := yaml.NewDecoder(f)
	var cfg Config
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return cfg
}

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
		},
		Commands: []*cli.Command{
			{
				Name:    "info",
				Aliases: []string{"i"},
				Usage:   "Displays info on given game server",
				Action: func(c *cli.Context) error {
					fmt.Printf("info %q", c.Args().Get(0))
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
	fmt.Println(".mima.yml")
	config := getConfig("mima.yml")

	fmt.Println(config)

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
