package webhooks

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"log"
)

//LogHandler is the handler for /log post endpoint. It deserializes the body according to LogPayload struct
// TODO - implement batching
// TODO - logging and adding more info to errors
func LogHandler() gin.HandlerFunc  {
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
		log.Println(lp)
		c.JSON(200,gin.H{
			"status":"ok",
		})
	}
}

func HealthzHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.JSON(200, gin.H{"status":"OK"}), // for json response
		c.String(200,"%s","OK") // for string response
	}
}