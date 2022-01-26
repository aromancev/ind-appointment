package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	URL           = "https://oap.ind.nl/oap/api/desks/AM/slots/?productKey=BIO&persons=2"
	telegramToken = "1797399084:AAEIltMLQEARSKTYAN6KUSkkKiyAog_cFYI"
	chatID        = 39818556
)

var (
	bot        *tgbotapi.BotAPI
	targetDate time.Time
)

func init() {
	var err error
	bot, err = tgbotapi.NewBotAPI(telegramToken)
	if err != nil {
		log.Fatal(err)
	}
	targetDate, err = time.Parse("2006-01-02", "2022-02-26")
	if err != nil {
		log.Fatal(err)
	}
}

func message(message string) {
	msg := tgbotapi.NewMessage(chatID, message)
	_, err := bot.Send(msg)
	if err != nil {
		log.Fatal(err)
	}
}

func receive() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatal(err)
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("[%s] (%d) %s", update.Message.From.UserName, update.Message.Chat.ID, update.Message.Text)
	}
}

type Response struct {
	Data []Appointment `json:"data"`
}

type Appointment struct {
	Date string `json:"date"`
	Time string `json:"endTime"`
}

func fetch() ([]string, error) {
	res, err := http.Get(URL)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var response Response
	err = json.Unmarshal(buf[6:], &response)
	if err != nil {
		return nil, err
	}
	var result []string
	for _, app := range response.Data {
		window, err := time.Parse("2006-01-02T15:04", app.Date+"T"+app.Time)
		if err != nil {
			return nil, err
		}
		if window.Before(targetDate) {
			result = append(result, fmt.Sprint(window))
		}
	}

	return result, nil
}

func main() {
	for {
		time.Sleep(3 * time.Second)
		results, err := fetch()
		if err != nil {
			message(err.Error())
			continue
		}

		if len(results) == 0 {
			continue
		}
		message(strings.Join(results, "\n"))
	}
}
