package consul

import (
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	"net/url"
)

func CreateAPIConfigs(ini *ini.File) (map[string]api.Config, error) {
	configs := map[string]api.Config{}
	for _, sec := range ini.Sections() {
		apiConfig := api.DefaultConfig()
		consulURL := sec.Key("url").String()

		u, err := url.Parse(consulURL)
		if err != nil {
			return configs, err
		}
		apiConfig.Address = u.Host
		apiConfig.Scheme = u.Scheme
		configs[sec.Name()] = *apiConfig
	}
	return configs, nil
}
