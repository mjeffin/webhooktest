package webhooks

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

//LogHandler is the handler for /log post endpoint. It deserializes the body according to LogPayload struct
// TODO - implement batching
// TODO - logging and adding more info to errors
func LogHandler(lpc chan LogPayload) gin.HandlerFunc  {
	return func(c *gin.Context) {
		j,err := io.ReadAll(c.Request.Body)
		if err != nil {
			log.Println(err)
		}
		var lp LogPayload
		err = json.Unmarshal(j,&lp)
		if err != nil {
			log.Println(err)
		}
		//log.Println(lp)
		lpc <- lp
		c.JSON(200,gin.H{
			"status":"ok",
		})
	}
}
// HealthzHandler - for the health endpoint. Since content type is not specified,sending response as text
func HealthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.JSON(200, gin.H{"status":"OK"}), // for json response
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
