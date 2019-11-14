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
	"time"
)

const (
	botMention = "@serverlessrandombot"
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
	} `json:"chat"`
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

// ProcessCommand takes in an input text (like "/dice 12") and returns what the bot should respond with
func ProcessCommand(text string, rng *rand.Rand) (string, error) {
	args := strings.Split(text, " ")
	cmd := strings.TrimSuffix(strings.ToLower(args[0]), botMention)

	switch cmd {
	case "/coin":
		if rng.Intn(2) == 1 {
			return "Heads", nil
		}
		return "Tails", nil
	case "/dice":
		num := 6
		var err error
		if len(args) >= 2 {
			num, err = strconv.Atoi(args[1])
		}

		if num <= 0 {
			return string(""), fmt.Errorf("dice should be > 0")
		}

		if err != nil {
			return fmt.Sprintf(
				"Couldn't parse number: %s",
				args[1],
			), nil
		}

		return fmt.Sprintf(
			"Rolled a %d",
			rng.Intn(num)+1,
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
			responses[rng.Intn(len(responses))],
			args[rng.Intn(len(args)-1)+1],
		), nil
	default:
		return string(""), fmt.Errorf("Command not found")
	}
}

// Dice responds to /api/dice
func Dice(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		return
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	var upd Update
	decoder.Decode(&upd)

	w.Write([]byte("K"))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	if upd.Message != nil && upd.Message.Text != "" && upd.Message.Text[0] == '/' {
		res, err := ProcessCommand(upd.Message.Text, rng)
		if err == nil {
			sendMessage(upd.Message.Chat.ID, res)
		}
	}
}
