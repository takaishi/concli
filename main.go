package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/hashicorp/consul/api"
	"net/url"
	"os"
)

type Config struct {
	URL string
}

func clientConfig() api.Config {
	var config Config
	_, err := toml.DecodeFile(fmt.Sprintf("%s/.cnodes", os.Getenv("HOME")), &config)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(config.URL)
	if err != nil {
		panic(err)
	}
	api_config := api.DefaultConfig()
	api_config.Address = u.Host
	api_config.Scheme = u.Scheme
	return *api_config
}

func main() {
	config := clientConfig()
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
