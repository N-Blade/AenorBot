package wordfilter

import (
	"bufio"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
)

var (
	badWordsDataURI = []string{
		"https://raw.githubusercontent.com/ilyankou/List-of-Dirty-Naughty-Obscene-and-Otherwise-Bad-Words/master/ru",
		"https://raw.githubusercontent.com/RobertJGabriel/Google-profanity-words/master/list.txt",
	}

	filter *Filter
)

type Filter struct {
	BadWords []string
}

func Init(dg *discordgo.Session) error {

	filter = &Filter{}

	for _, src := range badWordsDataURI {
		resp, err := http.Get(src)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			filter.BadWords = append(filter.BadWords, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return err
		}
	}

	dg.AddHandler(MessageFilter)

	return nil
}

func MessageFilter(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	msgWords := strings.Split(m.Content, " ")

	for _, msgWord := range msgWords {
		for _, word := range filter.BadWords {
			if strings.Contains(msgWord, word) {
				log.Infof("Bad word found: %s, pattern: %s", msgWord, word)
				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
