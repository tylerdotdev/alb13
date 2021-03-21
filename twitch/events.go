package twitch

import "github.com/alb13/twitch/eventsub"

const (
	SUB        = "sub"
	RESUB      = "resub"
	GIFTSUB    = "subgift"
	RAID       = "RAID"
	CHEER      = "cheer"
	MESSAGE    = "message"
	REDEMPTION = "redemption"
)

// Event represents an IRC event
type Event struct {
	Name              string
	SubEvent          SubEvent
	GiftSubEvent      GiftSubEvent
	ResubEvent        ResubEvent
	CheerEvent        CheerEvent
	RaidEvent         RaidEvent
	MessageEvent      MessageEvent
	NotificationEvent NotificationEvent
}

type SubEvent struct {
	User string `json:"user"`
	Tier string `json:"tier"`
}

type GiftSubEvent struct {
	User      string `json:"user"`
	Recipient string `json:"recipient"`
	Tier      string `json:"tier"`
}

type ResubEvent struct {
	User    string `json:"user"`
	Tier    string `json:"tier"`
	Months  string `json:"months"`
	Message string `json:"message"`
}

type CheerEvent struct {
	User        string `json:"user"`
	Bits        int    `json:"bits"`
	Message     string `json:"message"`
	IsAnonymous bool   `json:"isAnonymous"`
}

type RaidEvent struct {
	From    string `json:"from"`
	Viewers int    `json:"viewers"`
}

type MessageEvent struct {
	User    string `json:"user"`
	Message string `json:"message"`
}

// Reward represents redeemed channel points reward data
type Reward struct {
	ID     string `json:"id"`
	Title  string `json:"title"`
	Prompt string `json:"prompt"`
	Cost   int    `json:"cost"`
}

type NotificationEvent struct {
	ID                   string `json:"id"`
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
	BroadcasterUserName  string `json:"broadcaster_user_name"`
	UserID               string `json:"user_id"`
	UserLogin            string `json:"user_login"`
	UserName             string `json:"user_name"`
	Status               string `json:"status"`
	RedeemedAt           string `json:"redeemed_at"`
	Reward               Reward `json:"reward"`
}

// Notification represents a notification from EventSub
type Notification struct {
	Event        NotificationEvent     `json:"event"`
	Subscription eventsub.Subscription `json:"subscription"`
}
