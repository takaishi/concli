package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
)

func LoadConfig() (*ini.File, error) {
	f := fmt.Sprintf("%s/.cnodes", os.Getenv("HOME"))
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
