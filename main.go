package main

import (
	"github.com/takaishi/concli/cmd"
	"github.com/takaishi/concli/config"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "~/.concli",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "nodes",
			Usage: "List nodes",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "all",
				},
			},
			Action: func(c *cli.Context) error {

				configs, err := config.LoadConfig()
				if err != nil {
					return err
				}

				if c.Bool("all") {
					for _, cfg := range configs {
						cmd.PrintNodes(c, cfg)
					}
				} else {
					err = cmd.PrintNodes(c, configs[c.String("profile")])
					if err != nil {
						log.Fatalln(err)
					}
				}
				return cmd.PrintNodes(c, configs[c.String("profile")])
			},
		},
		{
			Name:  "services",
			Usage: "List services",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "all",
				},
				cli.StringFlag{
					Name:  "state",
					Value: "any",
				},
			},

			Action: func(c *cli.Context) error {

				configs, err := config.LoadConfig()
				if err != nil {
					return err
				}

				if c.Bool("all") {
					for _, cfg := range configs {
						cmd.PrintServices(c, cfg)
					}
				} else {
					err = cmd.PrintServices(c, configs[c.String("profile")])
					if err != nil {
						log.Fatalln(err)
					}
				}
				return cmd.PrintServices(c, configs[c.String("profile")])

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
