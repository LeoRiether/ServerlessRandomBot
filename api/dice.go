package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func sendMessage(chatID int64, text string) {
	token := os.Getenv("token")
	http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%d&text=%s", token, chatID, text))
}

func answerInlineQuery(chatID int64, results string) {

}

// User core.telegram.org/bots/api#user
type User struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
}

// Message https://core.telegram.org/bots/api#message
type Message struct {
	ID   int64  `json:"message_id"`
	Text string `json:"text"`
	From User   `json:"from"`
}

// InlineQuery https://core.telegram.org/bots/api#inlinequery
type InlineQuery struct {
	ID   int64 `json:"id"`
	From User  `json:"from"`
}

// Update https://core.telegram.org/bots/api#update
type Update struct {
	ID      int64        `json:"update_id"`
	Message *Message     `json:"message"`
	Inline  *InlineQuery `json:"inline_query"`
	Chat    struct {
		ID   int64  `json:"chat_id"`
		Type string `json:"type"`
	} `json:"chat"`
}

func handleCommand(chatID int64, args []string) {
	switch args[0] {
	case "/coin":
		if rand.Intn(2) == 1 {
			sendMessage(chatID, "Heads")
		} else {
			sendMessage(chatID, "Tails")
		}
	case "/dice":
		num := 6
		var err error
		if len(args) >= 2 {
			num, err = strconv.Atoi(args[1])
		}

		if err != nil {
			sendMessage(chatID, fmt.Sprintf(
				"Couln't parse number: %s",
				args[1],
			))
			break
		}

		sendMessage(chatID, fmt.Sprintf(
			"Number from 1 to %d:\nRolled a %d",
			num,
			rand.Intn(num)+1,
		))
	case "/list":
		if len(args) == 1 {
			sendMessage(chatID, fmt.Sprintf("Please input a space-separated list after the command, like so: /list rock paper scissors"))
			break
		}

		var responses = [...]string{"%s, clearly",
			"I choose %s",
			"Has to be %s",
			"%s, isn't it?",
			"It's %s",
			"%s is the chosen one",
			"Couldn't not be %s",
			"I declare %s to be victorious",
		}

		sendMessage(chatID, fmt.Sprintf(
			responses[rand.Intn(len(responses))],
			args[rand.Intn(len(args)-1)+1],
		))
	default:
		// do nothing
	}
}

// Dice responds to /api/dice with a random number between [0,4)
func Dice(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		return
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var upd Update
	decoder.Decode(&upd)

	w.Write([]byte("K"))

	if upd.Message != nil && upd.Message.Text != "" && upd.Message.Text[0] == '/' {
		handleCommand(upd.Chat.ID, strings.Split(upd.Message.Text, " "))
	}
}
