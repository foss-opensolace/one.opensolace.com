package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func New() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	if err := loadRouter(); err != nil {
		panic(err)
	}

	if err := loadDB(); err != nil {
		panic(err)
	}
}

func read(file string) []byte {
	data, err := os.ReadFile(filepath.Join("config", file))

	if err != nil {
		panic(err)
	}

	expanded := os.ExpandEnv(string(data))

	return []byte(expanded)
}
