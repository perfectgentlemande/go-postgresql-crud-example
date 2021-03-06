package main

import (
	"fmt"
	"os"

	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/api"
	"github.com/perfectgentlemande/go-postgresql-crud-example/internal/database/dbuser"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DBUser *dbuser.Config `yaml:"dbuser"`
	Server *api.Config    `yaml:"server"`
}

func readConfig(fileName string) (*Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	config := &Config{}
	err = yaml.NewDecoder(file).Decode(config)
	if err != nil {
		return nil, fmt.Errorf("failed to decode: %w", err)
	}

	return config, nil
}
