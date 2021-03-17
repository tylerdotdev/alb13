package server

import (
	"log"
	"strings"

	"antonlabs.io/alb13/twitch"
	"antonlabs.io/alb13/websocket"
)

func handleSub(pool *websocket.Pool, event twitch.Event) {
	message := websocket.BroadcastMessage{Event: event.Name, Data: event.SubEvent}
	pool.Broadcast <- message
}

func handleResub(pool *websocket.Pool, event twitch.Event) {
	message := websocket.BroadcastMessage{Event: event.Name, Data: event.ResubEvent}
	pool.Broadcast <- message
}

func handleGiftSub(pool *websocket.Pool, event twitch.Event) {
	message := websocket.BroadcastMessage{Event: event.Name, Data: event.GiftSubEvent}
	pool.Broadcast <- message
}

func handleCheer(pool *websocket.Pool, event twitch.Event) {
	if event.CheerEvent.Bits >= 100 {
		message := websocket.BroadcastMessage{Event: event.Name, Data: event.CheerEvent}
		pool.Broadcast <- message
	}

	if event.CheerEvent.Bits >= 500 && strings.Contains(event.CheerEvent.Message, "hydrate") {
		message := websocket.BroadcastMessage{Event: "hydrate"}
		pool.Broadcast <- message
	}
}

func handleRedemption(pool *websocket.Pool, event twitch.Event) {
	if event.NotificationEvent.Reward.Title == "Hydrate" {
		log.Println("Hydrate redeemed")
		message := websocket.BroadcastMessage{Event: "hydrate"}
		pool.Broadcast <- message
	}
}

func eventHandler(pool *websocket.Pool, event twitch.Event) {
	log.Println(event.Name)
	switch event.Name {
	case twitch.SUB:
		handleSub(pool, event)
		break
	case twitch.RESUB:
		handleResub(pool, event)
		break
	case twitch.GIFTSUB:
		handleGiftSub(pool, event)
		break
	case twitch.CHEER:
		handleCheer(pool, event)
		break
	case twitch.REDEMPTION:
		handleRedemption(pool, event)
	default:
		break
	}
}

func onTwitchEvent(pool *websocket.Pool, events chan twitch.Event) {
	for {
		event := <-events
		eventHandler(pool, event)
	}
}
