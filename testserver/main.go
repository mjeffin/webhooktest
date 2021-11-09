// testserver is a server to simulate the post endpoint url. not part of the main application
package main

import (
	"github.com/gin-gonic/gin"
	"log"
	"mjeffin/webhooks"
)

func main() {
	r := gin.Default()
	r.POST("/test",webhooks.TestPostEndpoint())
	err := r.Run(":8081") // listen and serve on 0.0.0.0:8081
	if err != nil {
		log.Fatalf("Error starting server")
	}
}

