package webhooks

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

//BatchRoutine is the background go routine which receives and stores the LogPayload struct from gin handler
//When a set batch size is reached OR the batch interval has passed the application should forward the collected records
//as an array to the post endpoint and clear the in-memory cache of objects.
func BatchRoutine(c chan LogPayload, batchSize int, batchInterval time.Duration)  {
	var payloads []LogPayload
	for {
		select {
		case p:= <-c  :
			payloads = append(payloads, p)
			fmt.Println("added new payload")
			fmt.Println(payloads)
			if len(payloads) == batchSize {
				err := processBatch(payloads);
				if err != nil {
					log.Fatalf("Unable to send the bactch to post endpoint. Quitting")
				}
				log.Println("batch processed successfully. Starting new batch")
				payloads = nil
			}
		}
	}
}

//processBatch sends the batch to api endpoint
// TODO - implement the function
func processBatch(payloads []LogPayload) error {
	return nil
}

func GetBatchConfig() (batchSize int, batchInterval time.Duration) {
	batchSize,err := strconv.Atoi(os.Getenv("BATCH_SIZE"))
	fmt.Println(batchSize)
	if err != nil {
		log.Fatalf("Error parsing batch size - %s\nexiting",err)
	}
	bis := os.Getenv("BATCH_INTERVAL")
	fmt.Println(bis)
	batchInterval, err = time.ParseDuration(bis)
	if err != nil {
		log.Fatalf("Error parsing batch interval - %v\nexiting",err)
	}
	return batchSize,batchInterval
}
