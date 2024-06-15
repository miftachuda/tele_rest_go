package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Define the Telegram API URL and the bot token
const telegramAPIURL = "https://api.telegram.org/bot%s/sendMessage"

// Struct for the request payload to Telegram
type TelegramMessage struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

// Struct for the incoming HTTP message
type IncomingMessage struct {
	Message string `json:"message"`
}

func main() {
	// Retrieve the bot token and chat ID from environment variables
	botToken := "7196089811:AAHIPG2vgOq3csjaLb2OT83kNDKkaPmWfvE"
	chatID := "643295256"
	http.HandleFunc("/forward", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Read the incoming request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Unmarshal the incoming JSON message
		var incomingMessage IncomingMessage
		if err := json.Unmarshal(body, &incomingMessage); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		// Create the Telegram message payload
		telegramMessage := TelegramMessage{
			ChatID: chatID,
			Text:   incomingMessage.Message,
		}

		// Marshal the payload to JSON
		payload, err := json.Marshal(telegramMessage)
		if err != nil {
			http.Error(w, "Failed to marshal JSON", http.StatusInternalServerError)
			return
		}

		// Send the message to the Telegram API
		url := fmt.Sprintf(telegramAPIURL, botToken)
		resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
		if err != nil || resp.StatusCode != http.StatusOK {
			http.Error(w, "Failed to send message to Telegram", http.StatusInternalServerError)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Message forwarded successfully"))
	})

	// Start the HTTP server on port 9999
	log.Println("Listening on port 7777")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
