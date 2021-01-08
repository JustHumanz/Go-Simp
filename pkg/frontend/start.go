package main

import (
	"github.com/JustHumanz/Go-simp/pkg/backend/utility/runfunc"
	discordhandler "github.com/JustHumanz/Go-simp/pkg/frontend/discord_handler"
	log "github.com/sirupsen/logrus"
)

func main() {
	discordhandler.StartInit("../../config.toml")
	log.Info("Frontend ready.......")

	runfunc.Run()
}
