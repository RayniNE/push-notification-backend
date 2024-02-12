package models

import "github.com/SherClockHolmes/webpush-go"

type Publisher struct {
	Name        string                  `json:"name,omitempty"`
	LastName    string                  `json:"last_name,omitempty"`
	Subscribers []*webpush.Subscription `json:"subscribers,omitempty"`
}

type SubscriberKeys struct {
	P256DH string `json:"p256dh,omitempty"`
	Auth   string `json:"auth,omitempty"`
}

type PublisherMessage struct {
	Message string `json:"message,omitempty"`
}

type SendMessageDTO struct {
	Index        int
	PubMessage   PublisherMessage
	Subscription *webpush.Subscription
}
