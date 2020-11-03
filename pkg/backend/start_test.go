package main

import (
	"testing"

	"github.com/JustHumanz/Go-simp/pkg/backend/cronjob"
	"github.com/JustHumanz/Go-simp/pkg/backend/runner"
	log "github.com/sirupsen/logrus"
)

func testStart(t *testing.T) {
	err := runner.StartInit(configfile)
	if err != nil {
		log.Fatal(err)
	}
	cronjob.InitCron()
	log.Info("Backend ready.......")

}
