package config

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type LeszmonitorConfigFile struct {
	GlobalParameters map[string]any            `yaml:"global_parameters"`
	Tags             map[string]map[string]any `yaml:"tags"`
	Monitors         map[string]map[string]any `yaml:"monitors"`
}

var (
	fileConfig       = &LeszmonitorConfigFile{}
	globalParameters = map[string]any{}
)

func LoadFile(path string) error {
	if path == "" {
		fileConfig = &LeszmonitorConfigFile{}
		globalParameters = map[string]any{}
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %w", path, err)
	}

	var cfg LeszmonitorConfigFile
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	if err := decoder.Decode(&cfg); err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("failed to parse config file %s: %w", path, err)
	}

	params := map[string]any{}
	flatten("", cfg.GlobalParameters, params)

	fileConfig = &cfg
	globalParameters = params
	return nil
}

func File() *LeszmonitorConfigFile {
	return fileConfig
}

func LookupGlobalParameter(key string) (any, bool) {
	value, ok := globalParameters[key]
	return value, ok
}

func flatten(prefix string, in map[string]any, out map[string]any) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		if nested, ok := v.(map[string]any); ok {
			flatten(key, nested, out)
			continue
		}

		out[key] = v
	}
}
