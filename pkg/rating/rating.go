package rating

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

const RatingChannelID = "997829476207571024"

func RatingUpdater(s *discordgo.Session) {

	ctx := context.Background()
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBAddress := os.Getenv("DB_ADDRESS")
	DBName := os.Getenv("DB_NAME")
	connectString := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", DBUser, DBPassword, DBAddress, DBName)
	db, err := sql.Open("mysql", connectString)
	if err != nil {
		log.Fatalf("Can't init rating updater, error: %s", err)
	}

	log.Info("Rating updater inited")
	queries := New(db)

	msgs, err := s.ChannelMessages(RatingChannelID, 1, "", "", "")
	if err != nil {
		log.Error(err)
	}

	var ratingMessage *discordgo.Message

	if len(msgs) != 0 {
		ratingMessage = msgs[0]
	} else {
		ratingMessage, err = s.ChannelMessageSend(RatingChannelID, "tmp")
		if err != nil {
			log.Error(err)
		}
	}

	for {
		ratingArr, err := queries.GetTop30ByRating(ctx)
		if err != nil {
			log.Errorf("Can't get top 30, error: %s", err)
		}

		fullMessage := "Ни один из игроков ещё не сыграл 10 матчей!"

		if len(ratingArr) != 0 {
			fullMessage = ""
			for i, player := range ratingArr {
				row := fmt.Sprintf("%d. %s, Уровень: %d, Рейтинг: %d, Количество матчей: %d\n",
					i+1, player.CharName, player.CharLevel+1, player.Rating, player.BattlesCount)
				fullMessage += row
			}
			s.ChannelMessageEdit(RatingChannelID, ratingMessage.ID, fullMessage)
		}

		time.Sleep(time.Second * 5)
	}
}
