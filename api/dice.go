package api

import (
	"encoding/json"
	"math/rand"
	"net/http"
)

type Response struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}

// Dice responds to /api/dice with a random number between [0,4)
func Dice(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		return
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	inp := make(map[string]interface{})
	decoder.Decode(&inp)

	response := Response{
		ChatID: int(inp["message"]["chat"]["id"]),
		Text: fmt.Sprintf("Rolled a %d", rand.Intn(4)),
	}

	out, _ := json.Marshal(response)
	w.Write(out)
}
