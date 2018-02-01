package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
)

func main() {
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
	fmt.Printf("%+v", nodes)
}
