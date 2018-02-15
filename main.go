package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	flags "github.com/jessevdk/go-flags"
	"io/ioutil"
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
		panic(err)
	}
	ini, err := ini.Load(data, f)
	if err != nil {
		panic(err)
	}
	for _, sec := range ini.Sections() {
		api_config := api.DefaultConfig()
		consul_url := sec.Key("url").String()

		u, err := url.Parse(consul_url)
		if err != nil {
			panic(err)
		}
		api_config.Address = u.Host
		api_config.Scheme = u.Scheme
		configs[sec.Name()] = *api_config
	}
	return configs, nil
}

func printNodes(config api.Config) {
	client, err := api.NewClient(&config)
	if err != nil {
		panic(err)
	}
	catalog := client.Catalog()
	query_options := api.QueryOptions{}
	nodes, _, err := catalog.Nodes(&query_options)
	if err != nil {
		panic(err)
	}

	health := client.Health()
	healthChecks, _, err := health.State("any", &query_options)
	if err != nil {
		panic(err)
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
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	configs, err := getConfigs()
	if err != nil {
		panic(err)
	}
	if opts.All {
		for _, cfg := range configs {
			printNodes(cfg)
		}
	}
	printNodes(configs[opts.Profile])
}
