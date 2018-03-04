package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"github.com/takaishi/concli/config"
	"github.com/takaishi/concli/consul"
	"log"
	"os"
	"github.com/fatih/color"
)

// Options is command line flags.
type Options struct {
	Profile string `short:"p" long:"profile" default:"DEFAULT" description:"DC section name to print."`
	All     bool   `short:"a" long:"all" description:"Print all nodes on all DC."`
	State string `short:"s" long:"state" default:"any" description:"Service state. Available any or passing, warning, critical."`
}

func printServices(config api.Config, state string) error {
	client, err := api.NewClient(&config)
	if err != nil {
		return errors.Wrap(err, "failed to create consul client")
	}
	queryOptions := api.QueryOptions{}

	health := client.Health()
	healthChecks, _, err := health.State(state, &queryOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch consul health states")
	}

	for _, healthCheck := range healthChecks {
		var stype string
		if healthCheck.Name == "Serf Health Status" {
			stype = "health"
		} else {
			stype = "service"
		}
		l := fmt.Sprintf("%-7s %-9s %-10s %s\n", stype, healthCheck.Status, healthCheck.Node, healthCheck.ServiceName)
		if healthCheck.Status == "passing" {
			color.Green(l)
		} else if healthCheck.Status == "warning" {
			color.Yellow(l)
		} else if healthCheck.Status == "critical" {
			color.Red(l)
		}

	}
	return nil
}

func main() {
	var opts Options

	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	ini, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	configs, err := consul.CreateAPIConfigs(ini)
	if err != nil {
		log.Fatalln(err)
	}

	if opts.All {
		for _, cfg := range configs {
			printServices(cfg, opts.State)
		}
	} else {
		err = printServices(configs[opts.Profile], opts.State)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
