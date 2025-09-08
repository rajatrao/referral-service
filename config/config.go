package config

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/config"
	"go.uber.org/fx"
)

const (
	configDir = "./config"
	baseFile  = "/base.yaml"
	// devFile = "/dev.yaml"
)

var Module = fx.Module(
	"config",
	fx.Provide(
		Load,
	),
)

func Load() (config.Provider, error) {
	var lookup config.LookupFunc = func(key string) (string, bool) {
		return os.LookupEnv(key)
	}
	expandOpts := config.Expand(lookup)
	cwd, err := filepath.Abs(configDir)
	if err != nil {
		return nil, fmt.Errorf("filepath abs %w", err)
	}

	fileOpts := config.File(cwd + baseFile)
	var ymlOpts []config.YAMLOption
	ymlOpts = append(ymlOpts, fileOpts)
	ymlOpts = append(ymlOpts, expandOpts)

	cfg, err := config.NewYAML(ymlOpts...)
	if err != nil {
		return nil, fmt.Errorf("config newyaml %w", err)
	}

	return cfg, nil
}
