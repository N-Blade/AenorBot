package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/N-Blade/AenorBot/pkg/twitch"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

func main() {

	dg, err := discordgo.New("Bot " + os.Getenv("DISCORD_AUTH_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	go twitch.StreamWatcher(dg)

	log.Info("AenorBot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}
