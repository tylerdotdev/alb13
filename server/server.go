package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"antonlabs.io/alb13/twitch"
	"antonlabs.io/alb13/twitch/eventsub"
	"antonlabs.io/alb13/twitch/irc"
	"antonlabs.io/alb13/websocket"
)

func handleChannelPoints(pool *websocket.Pool, event eventsub.Event) {
	if event.Reward.Title == "Hydrate" {
		log.Println("Hydrate redeemed")
		message := websocket.BroadcastMessage{Event: "hydrate"}
		pool.Broadcast <- message
	}
}

func handleNewNotification(pool *websocket.Pool, notification eventsub.Notification) {
	switch notification.Subscription.Type {
	case twitch.ChannelPointsType:
		handleChannelPoints(pool, notification.Event)
		break
	default:
		break
	}
}

func setupRoutes(pool *websocket.Pool) {
	http.HandleFunc("/notification", func(w http.ResponseWriter, r *http.Request) {
		handleNotifcation(pool, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleWs(pool, w, r)
	})
}

func setupTwitch(pool *websocket.Pool) {
	broadcasterID := os.Getenv("TWITCH_BROADCASTER_ID")
	channel := os.Getenv("TWITCH_CHANNEL")
	clientID := os.Getenv("TWITCH_CLIENT_ID")
	clientSecret := os.Getenv("TWITCH_CLIENT_SECRET")
	whCallbackURL := os.Getenv("WH_CALLBACK_URL")
	whSecret := os.Getenv("WH_SECRET")

	t := twitch.NewTwitch(broadcasterID, clientID, clientSecret, whCallbackURL, whSecret)
	t.SubscribeToChannelPointsRedemptions()
	t.StartIRC(channel, func(event irc.Event) {
		message := websocket.BroadcastMessage{Event: event.Type, Data: event.Data}
		pool.Broadcast <- message
		log.Println(message.Event)
	})
}

func Start() {
	pool := websocket.NewPool()
	go pool.Start()

	setupRoutes(pool)
	setupTwitch(pool)

	port := os.Getenv("PORT")

	fmt.Println("Starting server at port", ":"+port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
