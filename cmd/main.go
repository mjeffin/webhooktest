package main

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"mjeffin/webhooks"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	r := gin.Default()
	c := make(chan webhooks.LogPayload)
	go webhooks.BatchRoutine(c)
	r.GET("/healthz",webhooks.HealthzHandler())
	r.POST("/log",webhooks.LogHandler(c))
	log.Info("Started the webhook server")
	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatalf("Error starting server")
	}
}


