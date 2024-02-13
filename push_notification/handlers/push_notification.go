package handlers

import (
	"fmt"
	"net/http"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/gin-gonic/gin"
	"github.com/raynine/push-notification/models"
)

type PushNotificationHandler struct {
	Publisher       models.Publisher
	VAPIDPublicKey  string
	VAPIDPrivateKey string
}

// queue handles the messages that will be sent. It will handle one request at a time, giving it a semaphore like
// functionality
var queue = make(chan models.SendMessageDTO, 1)

// NewPushNotificationHandler created a PushNotificationHandler. It requires a VAPID public and private key.
func NewPushNotificationHandler(VAPIDPublicKey, VAPIDPrivateKey string) *PushNotificationHandler {
	return &PushNotificationHandler{
		VAPIDPublicKey:  VAPIDPublicKey,
		VAPIDPrivateKey: VAPIDPrivateKey,
		Publisher: models.Publisher{
			Name:        "Rayni",
			LastName:    "Nuñez Espino",
			Subscribers: make([]*webpush.Subscription, 0),
		},
	}
}

// Subscribe adds a webpush subscription to the publisher. It returns a 200 Status Code.
func (h *PushNotificationHandler) Subscribe(c *gin.Context) {
	var sub *webpush.Subscription

	err := c.BindJSON(&sub)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid subscriber",
		})
		return
	}

	h.Publisher.Subscribers = append(h.Publisher.Subscribers, sub)

	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, gin.H{
		"message": "You've subscribe successfully!",
	})
}

// Publish iterates over all the subscribers and it sends it to the queue. It returns a 200 Status Code.
func (h *PushNotificationHandler) Publish(c *gin.Context) {

	var pubMessage models.PublisherMessage

	if h.Publisher.Subscribers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No subscriber to publish",
		})
		return
	}

	err := c.BindJSON(&pubMessage)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	for index, sub := range h.Publisher.Subscribers {
		dto := models.SendMessageDTO{
			Index:        index,
			Subscription: sub,
			PubMessage:   pubMessage,
		}
		queue <- dto
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Messages are being processed!",
	})
}

// GetPublishers returns the current publisher
func (h *PushNotificationHandler) GetPublishers(c *gin.Context) {
	c.JSON(http.StatusOK, h.Publisher)
}

// GetSubscribers return all the subscriber of the publisher
func (h *PushNotificationHandler) GetSubscribers(c *gin.Context) {
	c.JSON(http.StatusOK, h.Publisher.Subscribers)
}

// SendNotifications is a goroutine that iterates over the queue buffered channel and sends the corresponding notification.
func (h *PushNotificationHandler) SendNotifications() {
	for sub := range queue {
		fmt.Printf("Received notification to queue, processing #%v...\n", sub.Index)
		// Send Notification
		resp, err := webpush.SendNotification([]byte(sub.PubMessage.Message), sub.Subscription, &webpush.Options{
			Subscriber:      "example@example.com",
			VAPIDPublicKey:  h.VAPIDPublicKey,
			VAPIDPrivateKey: h.VAPIDPrivateKey,
			TTL:             30,
		})

		if err != nil {
			fmt.Printf("An error ocurred sending notification: %v\n", err.Error())
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			fmt.Printf("An error ocurred sending notification: %v\n", err.Error())
			return
		}

		fmt.Printf("Finished processing notification #%v from queue, freeing queue...\n", sub.Index)
	}
}
