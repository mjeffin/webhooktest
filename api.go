package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io"
)

//LogHandler is the handler for /log post endpoint. It deserializes the body according to LogPayload struct
func LogHandler(lpc chan LogPayload) gin.HandlerFunc  {
	return func(c *gin.Context) {
		j,err := io.ReadAll(c.Request.Body)
		defer c.Request.Body.Close()
		if err != nil {
			log.Error("Error reading payload",err)
			c.JSON(500, gin.H{"status":err.Error()})
			return
		}
		var lp LogPayload
		err = json.Unmarshal(j,&lp)
		if err != nil {
			log.Error("Error unmarshalling json payload ",err)
			c.JSON(500, gin.H{"status":err.Error()})
			return
		}
		log.Info("Received and parsed webhook payload")
		go func() { // to avoid blocking even in the case BatchRoutine is not ready.
			lpc <- lp
		}()
		c.JSON(200,gin.H{
			"status":"ok",
		})
	}
}
// HealthzHandler - for the health endpoint. Since content type is not specified,sending response as text
func HealthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.JSON(200, gin.H{"status":"OK"}), // for json response
		log.Info("Received webhook on healthz endpoint")
		c.String(200,"%s","OK") // for string response
	}
}

//TestPostEndpoint is used for testig the test endpoint and is not part of the main application
func TestPostEndpoint() gin.HandlerFunc  {
	return func(c *gin.Context) {
		j,err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		fmt.Println(string(j))
		c.JSON(200,gin.H{
			"status":"ok",
		})
	}
}
