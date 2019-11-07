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

// type Response struct {
// 	ChatID int    `json:"chat_id"`
// 	Text   string `json:"text"`
// }

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

	w.Write([]byte("K"))

	// Please don't hack me
	response, _ := http.Get(fmt.Sprintf("https://api.telegram.org/bot806058245:AAFUjPUE0v8A8Ye5uNieNkkpyl87dRhH9ps/sendMessage?chat_id=%d&text=Rolled a %d", inp.Msg.Chatt.ID, rand.Intn(4)))
	log.Printf("Response: %v", response.StatusCode)
}
