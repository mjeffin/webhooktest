package main

import (
	"github.com/gin-gonic/gin"
	"mjeffin/webhooks"
)

func main() {
	r := gin.Default()
	c := make(chan webhooks.LogPayload)
	go webhooks.BatchRoutine(c)
	r.GET("/healthz",webhooks.HealthzHandler())
	r.POST("/log",webhooks.LogHandler(c))
	r.Run() // listen and serve on 0.0.0.0:8080
}


