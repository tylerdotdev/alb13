package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alb13/twitch"
	"github.com/alb13/websocket"
)

func setupRoutes(pool *websocket.Pool, t *twitch.Twitch) {
	http.HandleFunc("/notification", func(w http.ResponseWriter, r *http.Request) {
		handleNotifcation(t, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWs(pool, w, r)
	})
}

func setupTwitch(pool *websocket.Pool) twitch.Twitch {
	broadcasterID := os.Getenv("TWITCH_BROADCASTER_ID")
	channel := os.Getenv("TWITCH_CHANNEL")
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	whCallbackURL := os.Getenv("WH_CALLBACK_URL")
	whSecret := os.Getenv("WH_SECRET")

	t := twitch.NewTwitch(broadcasterID, clientID, clientSecret, channel, whCallbackURL, whSecret)
	t.SubscribeToChannelPointsRedemptions()
	go t.Connect()
	go onTwitchEvent(pool, t.Events)

	return t
}

func Start() {
	pool := websocket.NewPool()
	go pool.Start()

	t := setupTwitch(pool)
	setupRoutes(pool, &t)

	port := os.Getenv("PORT")

	fmt.Println("Starting server at port", ":"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
