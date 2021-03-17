package eventsub

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	BaseURL                 = "https://api.twitch.tv/helix"
	ChannelPointsRedemption = "channel.channel_points_custom_reward_redemption.add"
)

// Condition represents condition data for subscription
type Condition struct {
	BroadCasterID string `json:"broadcaster_user_id"`
	RewardID      string `json:"reward_id"`
}

// Transport represents transport data for subscription
type Transport struct {
	Method   string `json:"method"`
	Callback string `json:"callback"`
	Secret   string `json:"secret"`
}

// Subscription represents subscription data
type Subscription struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	Version   string    `json:"version"`
	Condition Condition `json:"condition"`
	Transport Transport `json:"transport"`
	CreatedAt string    `json:"created_at"`
	Cost      int       `json:"cost"`
}

// Subscriptions represents a list of Subscription types
type Subscriptions struct {
	Data []Subscription `json:"data"`
}

// Challenge represents challenge data from Twitch when sending a subscription request
type Challenge struct {
	Challenge    string       `json:"challenge"`
	Subscription Subscription `json:"subscription"`
}

// CreateSubscription creates a Twitch EventSub subscription
func CreateSubscription(subType string, broadcasterID string, clientID string, token string, whCallbackURL string, whSecret string) {
	subscription := Subscription{
		Type:    subType,
		Version: "1",
		Condition: Condition{
			BroadCasterID: broadcasterID,
		},
		Transport: Transport{
			Method:   "webhook",
			Callback: whCallbackURL + "/notification",
			Secret:   whSecret,
		},
	}

	postBody, err := json.Marshal(subscription)

	if err != nil {
		log.Println("Failed to marshal CreateSubscription payload", err)
	}

	req, err := http.NewRequest("POST", BaseURL+"/eventsub/subscriptions", bytes.NewBuffer(postBody))

	if err != nil {
		log.Println("Error reading CreateSubscription request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	if _, err := client.Do(req); err != nil {
		log.Println("Failed to send CreateSubscription request", err)
	}
}

// GetSubscriptions returns all current subscriptions
func GetSubscriptions(clientID string, token string) Subscriptions {
	req, err := http.NewRequest("GET", BaseURL+"/eventsub/subscriptions", nil)

	if err != nil {
		log.Println("Error reading GetSubscriptions request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	r, err := client.Do(req)

	if err != nil {
		log.Println("Failed to send GetSubscriptions request", err)
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Println("Failed to parse GetSubscriptions request body", err)
	}

	var subscriptions Subscriptions

	if err := json.Unmarshal(reqBody, &subscriptions); err != nil {
		log.Println("Error parsing GetSubscriptions json", err)
	}

	return subscriptions
}

// DeleteSubscription deletes a Twitch EventSub subscription
func DeleteSubscription(id string, clientID string, token string) {
	req, err := http.NewRequest("DELETE", BaseURL+"/eventsub/subscriptions?id="+id, nil)

	if err != nil {
		log.Println("Error reading DeleteSubscription request. ", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := http.Client{}

	if _, err := client.Do(req); err != nil {
		log.Println("Failed to send DeleteSubscription request", err)
	}
}
