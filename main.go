package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/olekukonko/tablewriter"
	"os"
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
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Node", "Address", "Dataenter"})
	for _, node := range nodes {
		table.Append([]string{node.Node, node.Address, node.Datacenter})
	}
	table.Render()

}
