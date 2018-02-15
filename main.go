package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	flags "github.com/jessevdk/go-flags"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/url"
	"os"
)

type Options struct {
	Profile string `short:"p" long:"profile" default:"DEFAULT"`
	All     bool   `short:"a" long:"all"`
}

func getConfigs() (map[string]api.Config, error) {
	configs := map[string]api.Config{}
	f := fmt.Sprintf("%s/.cnodes", os.Getenv("HOME"))
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return configs, errors.Wrapf(err, "failed to read %q", f)
	}
	ini, err := ini.Load(data, f)
	if err != nil {
		return configs, errors.Wrapf(err, "failed to load %q", f)
	}
	for _, sec := range ini.Sections() {
		api_config := api.DefaultConfig()
		consul_url := sec.Key("url").String()

		u, err := url.Parse(consul_url)
		if err != nil {
			return configs, errors.Wrapf(err, "failed to parse consul_url %q", consul_url)
		}
		api_config.Address = u.Host
		api_config.Scheme = u.Scheme
		configs[sec.Name()] = *api_config
	}
	return configs, nil
}

func printNodes(config api.Config) error {
	client, err := api.NewClient(&config)
	if err != nil {
		return errors.Wrap(err, "failed to create consul client")
	}
	catalog := client.Catalog()
	query_options := api.QueryOptions{}
	nodes, _, err := catalog.Nodes(&query_options)
	if err != nil {
		return errors.Wrapf(err, "failed to fetch consul nodes")
	}

	health := client.Health()
	healthChecks, _, err := health.State("any", &query_options)
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

func main() {
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		log.Fatalln(err)
	}

	configs, err := getConfigs()
	if err != nil {
		log.Fatalln(err)
	}
	if opts.All {
		for _, cfg := range configs {
			printNodes(cfg)
		}
	}
	err = printNodes(configs[opts.Profile])
	if err != nil {
		log.Fatalln(err)
	}
}
