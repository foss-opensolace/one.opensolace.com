package config

import "gopkg.in/yaml.v3"

var DB struct {
	PSQLConnectionString  string `yaml:"psql_connection_string"`
	MongoConnectionString string `yaml:"mongo_connection_string"`
}

func loadDB() error {
	if err := yaml.Unmarshal(read("db.yaml"), &DB); err != nil {
		return err
	}

	return nil
}
