package main

import (
	"testing"

	discordhandler "github.com/JustHumanz/Go-simp/pkg/frontend/discord_handler"
	log "github.com/sirupsen/logrus"
)

func testStart(t *testing.T) {
	err := discordhandler.StartInit("../../config.toml")
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Frontend ready.......")

}
