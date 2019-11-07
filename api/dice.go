package api

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
)

type Chat struct {
	ID int `json:"id"`
}

type Message struct {
	Chatt Chat `json:"chat"`
}

type Update struct {
	Msg Message `json:"message"`
}

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
	var inp Update
	decoder.Decode(&inp)

	log.Printf("Received %v\n", inp)

	response := Response{
		ChatID: inp.Msg.Chatt.ID,
		Text:   fmt.Sprintf("Rolled a %d", rand.Intn(4)),
	}

	out, _ := json.Marshal(response)
	w.Write(out)
}
