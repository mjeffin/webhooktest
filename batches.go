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
//batch size and interval can be fetched within the function itself, but I like passing configurations explicitly
// References - https://talks.golang.org/2012/concurrency.slide
func BatchRoutine(c chan LogPayload, batchSize int, batchInterval time.Duration)  {
	var payloads []LogPayload
	timeout := time.After(batchInterval)
	for {
		select {
		case p:= <-c:
			payloads = append(payloads, p)
			fmt.Printf("Current batch size - %d\n",len(payloads))
			if len(payloads) == batchSize {
				log.Println("Batch size exceeded. Resetting the batch")
				ResetBatch(payloads,c,batchSize,batchInterval)
				return
			}
		case <- timeout :
			log.Println("Batch interval elapsed. Resetting the batch. Current batch size -",len(payloads))
			ResetBatch(payloads,c,batchSize,batchInterval)
			return
		}
	}
}

//ResetBatch sends the current batch for processing
//Called by BatchRoutine when the reset conditions are met.
//Calls BatchRoutine as a go-routine and exits
//Exit the program if there is error processing the batch
//Could have done this using recursion inside BatchRoutine itself, but what's the base condition to stop infinite recursion ?
//Infinite recursion could have been stopped by calling itself as go routine and exiting immediately.
func ResetBatch(payloads []LogPayload, c chan LogPayload, batchSize int, batchInterval time.Duration) {
	err := processBatch(payloads)
	if err != nil {
		log.Fatalf("%v\nUnable to send the bactch to post endpoint.\n Quitting",err)
		os.Exit(1)
	}
	log.Println("batch processed successfully. Starting new batch")
	go BatchRoutine(c, batchSize,batchInterval)
}

//processBatch sends the batch to api endpoint
// TODO - implement the function
func processBatch(payloads []LogPayload) error {
	return nil
}

//GetBatchConfig process the environment variables and returns them in appropriate types
//BATCH_SIZE - integer number
//BATCH_INTERVAL - A duration string is a sequence of positive decimal numbers, each with optional fraction and a unit suffix,
//such as "300ms" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
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
