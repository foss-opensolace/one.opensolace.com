package main

import (
	"github.com/foss-opensolace/api.opensolace.com/internal/app"
	"github.com/foss-opensolace/api.opensolace.com/internal/config"
	"github.com/foss-opensolace/api.opensolace.com/internal/db"
)

func main() {
	config.New()
	db.New()
	app.New()
}
