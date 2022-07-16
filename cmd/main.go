package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/N-Blade/AenorBot/pkg/rating"
	"github.com/N-Blade/AenorBot/pkg/twitch"
	"github.com/N-Blade/AenorBot/pkg/wordfilter"

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

	err = wordfilter.Init(dg)
	if err != nil {
		log.Fatal("Can't init word filter", err)
	}
	log.Info("Word filter inited")

	go twitch.StreamWatcher(dg)
	go rating.RatingUpdater(dg)

	log.Info("AenorBot is now running.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
