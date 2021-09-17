package twitch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

type TwitchStreams struct {
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

func StreamWatcher(s *discordgo.Session) {
	Streams := make(map[string]*Stream)
	TwitchClientID := os.Getenv("TWITCH_CLIENT_ID")
	TwitchBearerToken := os.Getenv("TWITCH_BEARER_TOKEN")
	StreamsChannelID := os.Getenv("STREAMS_CHANNEL_ID")

	client := &http.Client{}

	for {
		req, err := http.NewRequest("GET", TwitchApiURL, nil)
		if err != nil {
			log.Error(err)
		}
		req.Header.Add("Client-ID", TwitchClientID)
		req.Header.Add("Authorization", "Bearer "+TwitchBearerToken)

		res, err := client.Do(req)
		if err != nil {
			log.Error(err)
		}

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Error(err)
		}

		res.Body.Close()

		var streams TwitchStreams
		if err := json.Unmarshal(body, &streams); err != nil {
			log.Error("Can unmarshal json", err)
		}

		if len(streams.Data) == 0 {
			for k, v := range Streams {
				err = s.ChannelMessageDelete(StreamsChannelID, v.MessageID)
				if err != nil {
					log.Error("Can't delete message with ID:", v, err)
				}
				delete(Streams, k)
			}
		}

		for _, stream := range streams.Data {
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
		for _, stream := range streams.Data {
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
