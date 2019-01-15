package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/takaishi/concli/consul"
	"io/ioutil"
	"os"
)

func LoadConfig() (map[string]api.Config, error) {
	ini, err := readConfig()
	if err != nil {
		return nil, err
	}
	configs, err := consul.CreateAPIConfigs(ini)
	if err != nil {
		return nil, err
	}

	return configs, nil
}

func readConfig() (*ini.File, error) {
	f := fmt.Sprintf("%s/.concli", os.Getenv("HOME"))
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read %q", f)
	}
	ini, err := ini.Load(data, f)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load %q", f)
	}
	return ini, nil
}
