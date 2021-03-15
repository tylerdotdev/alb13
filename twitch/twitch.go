package twitch

import (
	"antonlabs.io/alb13/twitch/auth"
	"antonlabs.io/alb13/twitch/eventsub"
	"antonlabs.io/alb13/twitch/irc"
)

const ChannelPointsType = "channel.channel_points_custom_reward_redemption.add"

type Twitch struct {
	AccessToken        auth.AccessToken
	BroadCasterID      string
	ClientID           string
	ClientSecret       string
	WebhookCallbackURL string
	WebhookSecret      string
}

// SetupTwitch setup up Twitch information
func NewTwitch(broadCasterID string, clientID string, clientSecret string, webhookCallbackURL string, webhookSecret string) Twitch {
	token := auth.GetAppAccessToken(clientID, clientSecret)

	return Twitch{
		AccessToken:        token,
		BroadCasterID:      broadCasterID,
		ClientID:           clientID,
		ClientSecret:       clientSecret,
		WebhookCallbackURL: webhookCallbackURL,
		WebhookSecret:      webhookSecret,
	}
}

func (twitch *Twitch) StartIRC(channel string, eventHandler func(event irc.Event)) {
	i := irc.NewIRC(channel)
	go i.Connect()
	go i.OnEvent(eventHandler)
}

func (twitch *Twitch) SubscribeToChannelPointsRedemptions() {
	subscriptions := eventsub.GetSubscriptions(twitch.ClientID, twitch.AccessToken.Token)

	for _, sub := range subscriptions.Data {
		if sub.Type == ChannelPointsType {
			eventsub.DeleteSubscription(sub.ID, twitch.ClientID, twitch.AccessToken.Token)
		}
	}

	eventsub.CreateSubscription(
		ChannelPointsType,
		twitch.BroadCasterID,
		twitch.ClientID,
		twitch.AccessToken.Token,
		twitch.WebhookCallbackURL,
		twitch.WebhookSecret,
	)
}
