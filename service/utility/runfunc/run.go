package runfunc

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func Run(Bot *discordgo.Session) {
	shutdown := make(chan int)
	//create a notification channel to shutdown
	sigChan := make(chan os.Signal, 1)

	//register for interupt (Ctrl+C) and SIGTERM (docker)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Info("Shutting down...")
		if Bot != nil {
			log.Info("Close bot wss session")
			err := Bot.Close()
			if err != nil {
				log.Error(err)
			}
			log.Info("Close database session")
			database.DbStop()
		}
		shutdown <- 1
	}()

	<-shutdown
}
