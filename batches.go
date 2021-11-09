package webhooks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

//BatchRoutine is the background go routine which receives and stores the LogPayload struct from gin handler
//When a set batch size is reached OR the batch interval has passed the application should forward the collected records
//as an array to the post endpoint and clear the in-memory cache of objects.
//clearing in-memory cache of objects is done by garbage collector as the go-routine exits after calling resetBatch
//Configs could have been passed from the main function too. That's how I implemented it initially
//Slightly faster that way, but slightly less readable and add tighter coupling too!
// References - https://talks.golang.org/2012/concurrency.slide
func BatchRoutine(c chan LogPayload)  {
	batchSize, batchInterval,postEndpoint := getBatchConfig()
	var payloads []LogPayload
	timeout := time.After(batchInterval)
	for {
		select {
		case p:= <-c:
			payloads = append(payloads, p)
			fmt.Printf("Current batch size - %d\n",len(payloads))
			if len(payloads) == batchSize {
				log.Println("Batch size exceeded. Resetting the batch")
				resetBatch(payloads,c,postEndpoint)
				return
			}
		case <- timeout :
			log.Println("Batch interval elapsed. Resetting the batch. Current batch size -",len(payloads))
			resetBatch(payloads,c,postEndpoint)
			return
		}
	}
}

//resetBatch sends the current batch for processing
//Called by BatchRoutine when the reset conditions are met.
//Calls BatchRoutine as a go-routine and exits
//Exit the program if there is error processing the batch
//ASSUMPTION - Post to test endpoint is done only if batch size > 0. why waste resource
//Could have done this using recursion inside BatchRoutine itself, but what's the base condition to stop infinite recursion ?
//Infinite recursion could have been stopped by calling itself as go routine and exiting immediately.
func resetBatch(payloads []LogPayload, c chan LogPayload, postEndpoint string) {
	go BatchRoutine(c) // immediately start the next batch
	currBatchSize := len(payloads)
	if currBatchSize >0 {
		err := processBatch(payloads,postEndpoint)
		if err != nil {
			log.Fatalf("%v\nUnable to send the bactch to post endpoint.\n Quitting",err)
			os.Exit(1)
		} else {
			log.Println("batch processed successfully. Starting new batch")
		}
	}
}

//processBatch sends the batch to api endpoint
func processBatch(payloads []LogPayload, postEndpoint string) error {
	b,err := json.Marshal(payloads)
	if err != nil {
		fmt.Println(err)
	}
	buf := bytes.NewBuffer(b)
	r,err := http.Post(postEndpoint,"application/json",buf)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(r.StatusCode)
	return nil
}

//GetBatchConfig process the environment variables and returns them in appropriate types
//BATCH_SIZE - integer number
//BATCH_INTERVAL - A duration string is a sequence of positive decimal numbers, each with optional fraction and a unit suffix,
//such as "300ms" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
//POST_ENDPOINT - A valid url
func getBatchConfig() (batchSize int, batchInterval time.Duration, postEndpoint string) {
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
	postEndpoint = os.Getenv("POST_ENDPOINT")
	_,err = url.ParseRequestURI(postEndpoint)
	if err != nil {
		log.Fatalf("Invalid post endpoint. Quitting")
	}
	return batchSize,batchInterval,postEndpoint
}
