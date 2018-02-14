package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	flags "github.com/jessevdk/go-flags"
	"io/ioutil"
	"net/url"
	"os"
)

type Config struct {
	URL string
}

type Options struct {
	Profile string `short:"p" long:"profile" default:"default"`
}

func clientConfig(profile string) (api.Config, error) {
	api_config := api.DefaultConfig()

	f := fmt.Sprintf("%s/.cnodes", os.Getenv("HOME"))
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return *api_config, err
	}
	ini, err := ini.Load(data, f)
	if err != nil {
		return *api_config, err
	}
	sec, err := ini.GetSection(profile)
	if err != nil {
		return *api_config, err
	}
	consul_url := sec.Key("url").String()
	fmt.Println(consul_url)

	u, err := url.Parse(consul_url)
	if err != nil {
		return *api_config, err
	}
	api_config.Address = u.Host
	api_config.Scheme = u.Scheme
	return *api_config, nil
}

func main() {
	var opts Options
	_, err := flags.Parse(&opts)

	if err != nil {
		os.Exit(1)
	}

	config, err := clientConfig(opts.Profile)
	if err != nil {
		panic(err)
	}

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

		fmt.Printf("%s\t%s\t%s\t%s\n", node.Node, node.Address, node.Datacenter, status)
	}
}
