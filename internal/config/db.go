package config

import "gopkg.in/yaml.v3"

var DB struct {
	ConnectionString string `yaml:"connection_string"`
}

func loadDB() error {
	if err := yaml.Unmarshal(read("db.yaml"), &DB); err != nil {
		return err
	}

	return nil
}
