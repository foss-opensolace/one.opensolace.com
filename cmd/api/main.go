package main

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/app"
	"github.com/foss-opensolace/api.opensolace.com/internal/config"
)

func main() {
	config.New()
	app.New()
}
