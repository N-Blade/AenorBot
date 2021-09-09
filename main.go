package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type TwitchStreamsJSON struct {
	Data []struct {
		ID           string      `json:"id"`
		UserID       string      `json:"user_id"`
		UserLogin    string      `json:"user_login"`
		UserName     string      `json:"user_name"`
		GameID       string      `json:"game_id"`
		GameName     string      `json:"game_name"`
		Type         string      `json:"type"`
		Title        string      `json:"title"`
		ViewerCount  int         `json:"viewer_count"`
		StartedAt    time.Time   `json:"started_at"`
		Language     string      `json:"language"`
		ThumbnailURL string      `json:"thumbnail_url"`
		TagIds       interface{} `json:"tag_ids"`
		IsMature     bool        `json:"is_mature"`
	} `json:"data"`
	Pagination struct {
	} `json:"pagination"`
}

type Stream struct {
	UserName     string
	Title        string
	ThumbnailURL string
	URL          string
	MessageID    string
}

const (
	TwitchApiURL = "https://api.twitch.tv/helix/streams?game_id=1905823964"
)

var (
	Streams           = make(map[string]*Stream)
	StreamsChannelID  = os.Getenv("STREAMS_CHANNEL_ID")
	DiscordAuthToken  = os.Getenv("DISCORD_AUTH_TOKEN")
	TwitchClientID    = os.Getenv("TWITCH_CLIENT_ID")
	TwitchBearerToken = os.Getenv("TWITCH_BEARER_TOKEN")
)

func main() {

	dg, err := discordgo.New("Bot " + DiscordAuthToken)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
	}

	go streamWatcher(dg)

	log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

func streamWatcher(s *discordgo.Session) {

	for {
		client := &http.Client{}
		req, err := http.NewRequest("GET", TwitchApiURL, nil)

		if err != nil {
			log.Error(err)
			return
		}
		req.Header.Add("Client-ID", TwitchClientID)
		req.Header.Add("Authorization", "Bearer "+TwitchBearerToken)

		res, err := client.Do(req)
		if err != nil {
			log.Error(err)
			return
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
			return
		}

		var streamsJSON TwitchStreamsJSON
		if err := json.Unmarshal(body, &streamsJSON); err != nil {
			log.Error("Can unmarshal json", err)
		}

		if len(streamsJSON.Data) == 0 {
			for k, v := range Streams {
				err = s.ChannelMessageDelete(StreamsChannelID, v.MessageID)
				if err != nil {
					log.Error("Can't delete message with ID:", v, err)
				}
				delete(Streams, k)
			}
		}

		for _, stream := range streamsJSON.Data {
			if _, ok := Streams[stream.ID]; ok {
				continue
			} else {
				fullURL := fmt.Sprintf("https://twitch.tv/%s", stream.UserLogin)
				log.Info("New stream detected!", fullURL)
				msg, err := s.ChannelMessageSend(StreamsChannelID,
					fmt.Sprintf("%s начал трансляцию! Ссылка на стрим: %s", stream.UserName, fullURL))
				if err != nil {
					log.Error("Can't send message", err)
				}

				Streams[stream.ID] = &Stream{
					UserName:     stream.UserName,
					Title:        stream.Title,
					ThumbnailURL: stream.ThumbnailURL,
					URL:          fullURL,
					MessageID:    msg.ID,
				}

			}
		}

		tempMap := make(map[string]*Stream)
		for _, stream := range streamsJSON.Data {
			tempMap[stream.ID] = &Stream{}
		}

		for k, v := range Streams {
			if _, ok := tempMap[k]; !ok {
				err = s.ChannelMessageDelete(StreamsChannelID, v.MessageID)
				if err != nil {
					log.Error("Can't delete message with ID:", v, err)
				}
				delete(Streams, k)
			}
		}

		time.Sleep(time.Second * 10)
	}

}
