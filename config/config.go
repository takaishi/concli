package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/url"
	"os"
)

func GetConfigs() (map[string]api.Config, error) {
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
		apiConfig := api.DefaultConfig()
		consulURL := sec.Key("url").String()

		u, err := url.Parse(consulURL)
		if err != nil {
			return configs, errors.Wrapf(err, "failed to parse consulURL %q", consulURL)
		}
		apiConfig.Address = u.Host
		apiConfig.Scheme = u.Scheme
		configs[sec.Name()] = *apiConfig
	}
	return configs, nil
}
