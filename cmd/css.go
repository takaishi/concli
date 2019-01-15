package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func PrintServices(c *cli.Context, cfg api.Config) error {
	client, err := api.NewClient(&cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create consul client")
	}
	queryOptions := api.QueryOptions{}

	health := client.Health()
	healthChecks, _, err := health.State(c.String("state"), &queryOptions)
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
