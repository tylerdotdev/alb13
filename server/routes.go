package server

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"antonlabs.io/alb13/twitch/eventsub"
	"antonlabs.io/alb13/websocket"
)

func verifySignature(signature string, id string, timestamp string, body []byte) bool {
	message := id + timestamp + string(body)
	h := hmac.New(sha256.New, []byte(os.Getenv("WH_SECRET")))
	h.Write([]byte(message))
	sha := "sha256=" + hex.EncodeToString(h.Sum(nil))
	return signature == sha
}

func handleWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		log.Println("Error upgrading connection", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client

	log.Println("Client connected to websocket")
}

func handleNotifcation(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Failed to parse request body", err)
	}

	if !verifySignature(
		r.Header.Get("Twitch-Eventsub-Message-Signature"),
		r.Header.Get("Twitch-Eventsub-Message-Id"),
		r.Header.Get("Twitch-Eventsub-Message-Timestamp"),
		body,
	) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "webhook_callback_verification" {
		var challenge eventsub.Challenge

		if err := json.Unmarshal(body, &challenge); err != nil {
			log.Println("Error parsing challenge json", err)
		}

		w.Write([]byte(challenge.Challenge))
	}

	if r.Header.Get("Twitch-Eventsub-Message-Type") == "notification" {
		var notification eventsub.Notification

		if err := json.Unmarshal(body, &notification); err != nil {
			log.Println("Error parsing challenge json", err)
		}

		w.WriteHeader(http.StatusNoContent)
		handleNewNotification(pool, notification)
	}
}
