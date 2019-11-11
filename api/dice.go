package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func sendMessage(chatID int64, text string) {
	token := os.Getenv("TOKEN")
	text = url.QueryEscape(text)
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
	Chat struct {
		ID   int64  `json:"id"`
		Type string `json:"type"`
	}
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
}

func handleCommand(text string) (string, error) {
	args := strings.Split(text, " ")
	cmd := strings.TrimSuffix(strings.ToLower(args[0]), "@serverlessrandombot")

	switch cmd {
	case "/coin":
		if rand.Intn(2) == 1 {
			return "Heads", nil
		} else {
			return "Tails", nil
		}
	case "/dice":
		num := 6
		var err error
		if len(args) >= 2 {
			num, err = strconv.Atoi(args[1])
		}

		if err != nil {
			return fmt.Sprintf(
				"Couln't parse number: %s",
				args[1],
			), nil
		}

		return fmt.Sprintf(
			"Rolled a %d",
			rand.Intn(num)+1,
		), nil
	case "/list":
		if len(args) == 1 {
			// sendMessage(chatID, fmt.Sprintf("Please input a space-separated list after the command, like so: /list rock paper scissors"))
			return "Segmentation fault", nil
		}

		var responses = [...]string{
			"%s, clearly",
			"I choose %s",
			"Has to be %s",
			"%s, isn't it?",
			"It's %s",
			"%s is the chosen one",
			"Couldn't not be %s",
			"I declare %s to be victorious",
		}

		return fmt.Sprintf(
			responses[rand.Intn(len(responses))],
			args[rand.Intn(len(args)-1)+1],
		), nil
	default:
		return string(""), fmt.Errorf("Command not found")
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
		res, err := handleCommand(upd.Message.Text)
		if err == nil {
			sendMessage(upd.Message.Chat.ID, res)
		}
	}
}
