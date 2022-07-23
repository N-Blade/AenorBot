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

	msgs, err := s.ChannelMessages(RatingChannelID, 2, "", "", "")
	if err != nil {
		log.Error(err)
	}

	var soloRatingMessage, partyRatingMessage *discordgo.Message

	if len(msgs) == 0 {
		soloRatingMessage, err = s.ChannelMessageSend(RatingChannelID, "Пока здесь ничего нет, но скоро будет!")
		if err != nil {
			log.Error(err)
		}
		partyRatingMessage, err = s.ChannelMessageSend(RatingChannelID, "Пока здесь ничего нет, но скоро будет!")
		if err != nil {
			log.Error(err)
		}
	} else {
		soloRatingMessage = msgs[1]
		partyRatingMessage = msgs[0]
	}

	for {
		soloRatingArr, err := queries.GetTop30SoloRating(ctx)
		if err != nil {
			log.Errorf("Can't get top 30 solo, %s", err)
		}

		soloFullMessage := "Одиночный рейтинг\nНи один из игроков ещё не сыграл 10 одиночных матчей!"

		if len(soloRatingArr) != 0 {
			soloFullMessage = "Одиночный рейтинг\n"
			for i, player := range soloRatingArr {
				row := fmt.Sprintf("%d. %s (%s), Уровень: %d, Рейтинг: %d, Количество матчей: %d\n",
					i+1, player.CharName, player.CharGuildName, player.CharLevel+1, player.SoloRating, player.MatchCount)
				soloFullMessage += row
			}
			_, err = s.ChannelMessageEdit(RatingChannelID, soloRatingMessage.ID, soloFullMessage)
			if err != nil {
				log.Error(err)
			}
		}

		partyFullMessage := "Групповой рейтинг\nНи один из игроков ещё не сыграл 10 групповых матчей!"

		partyRatingArr, err := queries.GetTop30PartyRating(ctx)
		if err != nil {
			log.Errorf("Can't get top 30 party, %s", err)
		}

		if len(partyRatingArr) != 0 {
			partyFullMessage = "Групповой рейтинг\n"
			for i, player := range partyRatingArr {
				row := fmt.Sprintf("%d. %s (%s), Уровень: %d, Рейтинг: %d, Количество матчей: %d\n",
					i+1, player.CharName, player.CharGuildName, player.CharLevel+1, player.PartyRating, player.MatchCount)
				partyFullMessage += row
			}
			_, err = s.ChannelMessageEdit(RatingChannelID, partyRatingMessage.ID, partyFullMessage)
			if err != nil {
				log.Error(err)
			}
		}

		time.Sleep(time.Second * 5)
	}
}
