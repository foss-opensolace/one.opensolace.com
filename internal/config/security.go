package config

import "gopkg.in/yaml.v3"

var Security struct {
	JWTSecret string `yaml:"jwt_secret"`
}

func loadSecurity() error {
	if err := yaml.Unmarshal(read("security.yaml"), &Security); err != nil {
		return err
	}

	return nil
}
