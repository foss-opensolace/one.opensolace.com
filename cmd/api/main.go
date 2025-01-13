package main

import (
	"github.com/foss-opensolace/one.opensolace.com/internal/api/service"
	"github.com/foss-opensolace/one.opensolace.com/internal/app"
	"github.com/foss-opensolace/one.opensolace.com/internal/config"
	"github.com/foss-opensolace/one.opensolace.com/internal/db"
)

func main() {
	config.New()
	db.New()
	service.New()
	app.New()
}
