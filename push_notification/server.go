package push_notification

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/raynine/push-notification/push_notification/handlers"
)

type Server struct {
	VAPIDPublicKey  string
	VAPIDPrivateKey string
}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) GenerateVAPIDKeys() {
	s.VAPIDPrivateKey = "w-SBDglSOKUR42opRhEnTn8xtei1nmq1RslTDpiGv98"
	s.VAPIDPublicKey = "BNTyUF41b8UXl7htTMPUyH3VJ_bQemJvK2OlSmFT9cLqzCCeVYS5lDPpmQhoze4YLyen6uAG-JCaRyLY66OK5jk"
}

func (s *Server) Init() {
	r := gin.Default()
	r.Use(cors.Default())

	handler := handlers.NewPushNotificationHandler(s.VAPIDPublicKey, s.VAPIDPrivateKey)
	go handler.SendNotifications()

	r.GET("/publishers", handler.GetPublishers)
	r.GET("/subscribers", handler.GetSubscribers)
	r.POST("/subscribe", handler.Subscribe)
	r.POST("/publish", handler.Publish)

	panic(http.ListenAndServe(":8080", r))

}
