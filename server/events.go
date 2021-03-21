package server

import (
	"log"
	"strings"

	"github.com/alb13/twitch"
	"github.com/alb13/websocket"
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

func handleRaid(pool *websocket.Pool, event twitch.Event) {
	if event.RaidEvent.Viewers >= 25 {
		message := websocket.BroadcastMessage{Event: event.Name, Data: event.RaidEvent}
		pool.Broadcast <- message
	}
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

func handleMessage(pool *websocket.Pool, event twitch.Event) {
	message := websocket.BroadcastMessage{Event: event.Name, Data: event.MessageEvent}
	pool.Broadcast <- message
}

func handleRedemption(pool *websocket.Pool, event twitch.Event) {
	if event.NotificationEvent.Reward.Title == "Hydrate" {
		log.Println("Hydrate redeemed")
		message := websocket.BroadcastMessage{Event: "hydrate"}
		pool.Broadcast <- message
	}
}

func eventHandler(pool *websocket.Pool, event twitch.Event) {
	if event.Name != twitch.MESSAGE {
		log.Println(event.Name)
	}

	switch event.Name {
	case twitch.SUB:
		handleSub(pool, event)
	case twitch.RESUB:
		handleResub(pool, event)
	case twitch.GIFTSUB:
		handleGiftSub(pool, event)
	case twitch.RAID:
		handleRaid(pool, event)
	case twitch.CHEER:
		handleCheer(pool, event)
	case twitch.MESSAGE:
		handleMessage(pool, event)
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
