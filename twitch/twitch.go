package twitch

import (
	"log"

	"github.com/alb13/twitch/auth"
	"github.com/alb13/twitch/eventsub"
	irc "github.com/gempir/go-twitch-irc/v2"
)

type Config struct {
	AccessToken        auth.AccessToken
	BroadcasterID      string
	ClientID           string
	ClientSecret       string
	IRCChannel         string
	WebhookCallbackURL string
	WebhookSecret      string
}

type Twitch struct {
	Config Config
	Events chan Event
}

// SetupTwitch setup up Twitch information
func NewTwitch(broadCasterID string, clientID string, clientSecret string, ircChannel string, webhookCallbackURL string, webhookSecret string) Twitch {
	token := auth.GetAppAccessToken(clientID, clientSecret)

	config := Config{
		AccessToken:        token,
		BroadcasterID:      broadCasterID,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		IRCChannel:         ircChannel,
		WebhookCallbackURL: webhookCallbackURL,
		WebhookSecret:      webhookSecret,
	}

	return Twitch{
		Config: config,
		Events: make(chan Event),
	}
}

func (twitch *Twitch) Connect() {
	client := irc.NewAnonymousClient()

	client.Join(twitch.Config.IRCChannel)

	go client.OnConnect(func() {
		log.Println("Connected to", twitch.Config.IRCChannel)
	})

	go client.OnUserNoticeMessage(twitch.handleUserNotice)
	go client.OnPrivateMessage(twitch.handlePrivateMessage)

	err := client.Connect()
	if err != nil {
		log.Println("Failed to connect to IRC:", err)
	}
}

func (twitch *Twitch) SubscribeToChannelPointsRedemptions() {
	subscriptions := eventsub.GetSubscriptions(twitch.Config.ClientID, twitch.Config.AccessToken.Token)

	for _, sub := range subscriptions.Data {
		if sub.Type == eventsub.ChannelPointsRedemption {
			eventsub.DeleteSubscription(sub.ID, twitch.Config.ClientID, twitch.Config.AccessToken.Token)
		}
	}

	eventsub.CreateSubscription(
		eventsub.ChannelPointsRedemption,
		twitch.Config.BroadcasterID,
		twitch.Config.ClientID,
		twitch.Config.AccessToken.Token,
		twitch.Config.WebhookCallbackURL,
		twitch.Config.WebhookSecret,
	)
}

func (twitch *Twitch) HandleEventSubNotification(notification Notification) {
	switch notification.Subscription.Type {
	case eventsub.ChannelPointsRedemption:
		event := Event{Name: REDEMPTION, NotificationEvent: notification.Event}
		twitch.Events <- event
		break
	default:
		break
	}
}

func (twitch *Twitch) handleSub(message irc.UserNoticeMessage) {
	sub := SubEvent{
		User: message.User.DisplayName,
		Tier: message.MsgParams["msg-param-sub-plan"],
	}

	event := Event{Name: SUB, SubEvent: sub}
	twitch.Events <- event
}

func (twitch *Twitch) handleResub(message irc.UserNoticeMessage) {
	resub := ResubEvent{
		User:    message.User.DisplayName,
		Tier:    message.MsgParams["msg-param-sub-plan"],
		Months:  message.MsgParams["msg-param-cumulative-months"],
		Message: message.Message,
	}

	event := Event{Name: RESUB, ResubEvent: resub}
	twitch.Events <- event
}

func (twitch *Twitch) handleGiftSub(message irc.UserNoticeMessage) {
	subGift := GiftSubEvent{
		User:      message.User.DisplayName,
		Recipient: message.MsgParams["msg-param-recipient-display-name"],
		Tier:      message.MsgParams["msg-param-sub-plan"],
	}

	event := Event{Name: GIFTSUB, GiftSubEvent: subGift}
	twitch.Events <- event
}

func (twitch *Twitch) handleCheer(message irc.PrivateMessage) {
	cheer := CheerEvent{
		User:    message.User.DisplayName,
		Message: message.Message,
		Bits:    message.Bits,
	}

	event := Event{Name: CHEER, CheerEvent: cheer}
	twitch.Events <- event
}

func (twitch *Twitch) handleUserNotice(message irc.UserNoticeMessage) {
	switch message.MsgID {
	case SUB:
		twitch.handleSub(message)
		break
	case RESUB:
		twitch.handleResub(message)
		break
	case GIFTSUB:
		twitch.handleGiftSub(message)
		break
	default:
		break
	}
}

func (twitch *Twitch) handlePrivateMessage(message irc.PrivateMessage) {
	if message.Bits > 0 {
		twitch.handleCheer(message)
	}
}
