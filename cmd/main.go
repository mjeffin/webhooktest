// executable for the application
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
	//logging is already enabled as gin is being run in debug mode by default
	r := gin.Default()
	c := make(chan webhooks.LogPayload)
	go webhooks.BatchRoutine(c)
	r.POST("/log",webhooks.LogHandler(c))
	r.GET("/healthz",webhooks.HealthzHandler())
	log.Println("starting webhook server") // we can't print easily after starting server as Run will block forever
	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Fatalf("Error starting server")
	}

}


