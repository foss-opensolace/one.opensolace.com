package config

import "gopkg.in/yaml.v3"

var Router struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

func loadRouter() error {
	if err := yaml.Unmarshal(read("router.yaml"), &Router); err != nil {
		return err
	}

	return nil
}
