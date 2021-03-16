package irc

import (
	"log"

	"github.com/gempir/go-twitch-irc/v2"
)

const (
	SUB     = "sub"
	RESUB   = "resub"
	GIFTSUB = "subgift"
	CHEER   = "cheer"
)

type IRC struct {
	Client  *twitch.Client
	Channel string
	Events  chan Event
}

// NewIRC creates a new Twitch IRC client
func NewIRC(channel string) IRC {
	client := twitch.NewAnonymousClient()

	irc := IRC{
		Client:  client,
		Channel: channel,
		Events:  make(chan Event),
	}

	return irc
}

// Connect connects to configured Twitch channel
func (irc *IRC) Connect() {
	irc.Client.Join(irc.Channel)

	go irc.Client.OnConnect(func() {
		log.Println("Connected to:", irc.Channel)
	})

	go irc.Client.OnUserNoticeMessage(irc.handleUserNotice)
	go irc.Client.OnPrivateMessage(irc.handlePrivateMessage)

	err := irc.Client.Connect()
	if err != nil {
		log.Println("Failed to connect to IRC:", err)
	}
}

func (irc *IRC) OnEvent(callback func(event Event)) {
	for {
		event := <-irc.Events
		callback(event)
	}
}

func (irc *IRC) handleSub(message twitch.UserNoticeMessage) {
	sub := SubEvent{
		User: message.User.DisplayName,
		Tier: message.MsgParams["msg-param-sub-plan"],
	}

	event := Event{Type: SUB, Data: sub}
	irc.Events <- event
}

func (irc *IRC) handleResub(message twitch.UserNoticeMessage) {
	resub := ResubEvent{
		User:    message.User.DisplayName,
		Tier:    message.MsgParams["msg-param-sub-plan"],
		Months:  message.MsgParams["msg-param-cumulative-months"],
		Message: message.Message,
	}

	event := Event{Type: RESUB, Data: resub}
	irc.Events <- event
}

func (irc *IRC) handleGiftSub(message twitch.UserNoticeMessage) {
	subGift := GiftSubEvent{
		User:      message.User.DisplayName,
		Recipient: message.MsgParams["msg-param-recipient-display-name"],
		Tier:      message.MsgParams["msg-param-sub-plan"],
	}

	event := Event{Type: GIFTSUB, Data: subGift}
	irc.Events <- event
}

func (irc *IRC) handleCheer(message twitch.PrivateMessage) {
	cheer := CheerEvent{
		User:    message.User.DisplayName,
		Message: message.Message,
		Bits:    message.Bits,
	}

	event := Event{Type: CHEER, Data: cheer}
	irc.Events <- event
}

func (irc *IRC) handleUserNotice(message twitch.UserNoticeMessage) {
	switch message.MsgID {
	case SUB:
		irc.handleSub(message)
		break
	case RESUB:
		irc.handleResub(message)
		break
	case GIFTSUB:
		irc.handleGiftSub(message)
		break
	default:
		break
	}
}

func (irc *IRC) handlePrivateMessage(message twitch.PrivateMessage) {
	if message.Bits > 0 {
		irc.handleCheer(message)
	}
}
