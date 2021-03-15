package irc

// Event represents an IRC event
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
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
	FromUser string `json:"fromUser"`
	Viewers  int    `json:"viewers"`
}
