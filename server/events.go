package server

import (
	"log"
	"strings"

	"antonlabs.io/alb13/twitch"
	"antonlabs.io/alb13/websocket"
)

func eventHandler(pool *websocket.Pool, event twitch.Event) {
	log.Println(event.Name)
	switch event.Name {
	case twitch.SUB:
		message := websocket.BroadcastMessage{Event: event.Name, Data: event.SubEvent}
		pool.Broadcast <- message
		break
	case twitch.RESUB:
		message := websocket.BroadcastMessage{Event: event.Name, Data: event.ResubEvent}
		pool.Broadcast <- message
		break
	case twitch.GIFTSUB:
		message := websocket.BroadcastMessage{Event: event.Name, Data: event.GiftSubEvent}
		pool.Broadcast <- message
		break
	case twitch.CHEER:
		if event.CheerEvent.Bits >= 100 {
			message := websocket.BroadcastMessage{Event: event.Name, Data: event.CheerEvent}
			pool.Broadcast <- message
		}

		if event.CheerEvent.Bits >= 500 && strings.Contains(event.CheerEvent.Message, "hydrate") {
			message := websocket.BroadcastMessage{Event: "hydrate"}
			pool.Broadcast <- message
		}
		break
	case twitch.REDEMPTION:
		if event.NotificationEvent.Reward.Title == "Hydrate" {
			log.Println("Hydrate redeemed")
			message := websocket.BroadcastMessage{Event: "hydrate"}
			pool.Broadcast <- message
		}
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
