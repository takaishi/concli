package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

func PrintNodes(c *cli.Context, cfg api.Config) error {
	client, err := api.NewClient(&cfg)
	if err != nil {
		return errors.Wrap(err, "failed to create consul client")
	}
	catalog := client.Catalog()
	queryOptions := api.QueryOptions{}
	nodes, _, err := catalog.Nodes(&queryOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch consul nodes")
	}

	health := client.Health()
	healthChecks, _, err := health.State("any", &queryOptions)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch consul health states")
	}

	for _, node := range nodes {
		status := "passing"
		for _, healthCheck := range healthChecks {
			if healthCheck.Node == node.Node && healthCheck.Status != "passing" {
				status = healthCheck.Status
			}
		}

		l := fmt.Sprintf("%-9s %-16s %s %s\n", status, node.Address, node.Node, node.Datacenter)
		if status == "passing" {
			color.Green(l)
		} else if status == "critical" {
			color.Red(l)

		}
	}
	return nil
}
