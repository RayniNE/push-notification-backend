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

func (h *PushNotificationHandler) Publish(c *gin.Context) {
	if h.Publisher.Subscribers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "No subscriber to publish",
		})
		return
	}

	for _, sub := range h.Publisher.Subscribers {
		// Send Notification
		resp, err := webpush.SendNotification([]byte("Notification in real time!"), sub, &webpush.Options{
			Subscriber:      "example@example.com",
			VAPIDPublicKey:  h.VAPIDPublicKey,
			VAPIDPrivateKey: h.VAPIDPrivateKey,
			TTL:             30,
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "An error happened doing the http call: " + err.Error(),
			})
			return
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Error sending messages to subscribers, got status_code: %v", resp.StatusCode),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Message published to all subscribers!",
	})
}

func (h *PushNotificationHandler) GetPublishers(c *gin.Context) {
	c.JSON(http.StatusOK, h.Publisher)
}

func (h *PushNotificationHandler) GetSubscribers(c *gin.Context) {
	c.JSON(http.StatusOK, h.Publisher.Subscribers)
}
