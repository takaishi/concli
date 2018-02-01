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

func main() {
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
	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		panic(err)
	}
	catalog := client.Catalog()
	query_options := api.QueryOptions{}
	nodes, _, err := catalog.Nodes(&query_options)
	if err != nil {
		panic(err)
	}

	for _, node := range nodes {
		fmt.Printf("%s\t%s\t%s\n", node.Node, node.Address, node.Datacenter)
	}
}
