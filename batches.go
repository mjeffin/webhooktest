package webhooks

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	maxReties = 3
	secondsToWait = 2
)

//BatchRoutine is the background go routine which receives and stores the LogPayload struct from gin handler
//When a set batch size is reached OR the batch interval has passed the application should forward the collected records
//as an array to the post endpoint and clear the in-memory cache of objects.
//clearing in-memory cache of objects is done by garbage collector as the go-routine exits after calling processBatch
//Configs could have been passed from the main function too. That's how I implemented it initially
//Slightly faster that way, but slightly less readable and add tighter coupling too!
// References - https://talks.golang.org/2012/concurrency.slide
func BatchRoutine(c chan LogPayload)  {
	batchSize, batchInterval,postEndpoint := getBatchConfig()
	log.WithFields(log.Fields{"batch_size":batchSize,"batch_interval":batchInterval.String()}).Info("starting new batch in go-routine")
	var payloads []LogPayload
	timeout := time.After(batchInterval)
	for {
		select {
		case p:= <-c:
			payloads = append(payloads, p)
			log.WithFields(log.Fields{"current_batch_size":len(payloads)}).Info("payload added to batch")
			if len(payloads) == batchSize {
				log.Info("Batch size exceeded. Processing the batch")
				processBatch(payloads,c,postEndpoint)
				return
			}
		case <- timeout :
			log.Info("Batch interval elapsed. Processing the batch")
			processBatch(payloads,c,postEndpoint)
			return
		}
	}
}

//processBatch sends the current batch for processing
//Called by BatchRoutine when the reset conditions are met.
//Calls BatchRoutine as a go-routine and exits
//Exit the program if there is error processing the batch
//ASSUMPTION - Post to test endpoint is done only if batch size > 0. why waste resource
//Could have done this using recursion inside BatchRoutine itself, but what's the base condition to stop infinite recursion ?
//Infinite recursion could have been stopped by calling itself as go routine and exiting immediately.
func processBatch(payloads []LogPayload, c chan LogPayload, postEndpoint string) {
	go BatchRoutine(c) // immediately start the next batch
	currBatchSize := len(payloads)
	log.WithFields(log.Fields{"batch_size":currBatchSize}).Info("Received a new batch to process")
	if currBatchSize < 1 {
		return
	}
	for i:=0;i<maxReties;i++ {
		err := sendBatch(payloads,postEndpoint)
		if err != nil {
			if i == (maxReties -1) {
				log.WithFields(log.Fields{"retry_counter":maxReties}).Fatalf("Unable to send the batch to post endpoint after maximum reties. Quitting")
			} else {
				log.WithFields(log.Fields{"retry_counter":i+1}).Errorf("Sending batch to post endpoint failed. Retrying after %d seconds",secondsToWait)
				time.Sleep(secondsToWait * time.Second)
			}
		} else {
			return
		}
	}
}

//sendBatch sends the batch to api endpoint
func sendBatch(payloads []LogPayload, postEndpoint string) error {
	b,err := json.Marshal(payloads)
	if err != nil {
		log.WithFields(log.Fields{"error":err}).Error()
		return errors.New("error marshalling records to send to post endpoint")
	}
	buf := bytes.NewBuffer(b)
	startTime := time.Now()
	r,err := http.Post(postEndpoint,"application/json",buf)
	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	if err != nil {
		log.WithFields(log.Fields{"error":err,"post_duration": elapsed.String(),
			"batch_size": len(payloads)}).Error("error sending records to post endpoint")
		return errors.New("error sending records to post endpoint")
	}
	if !(r.StatusCode >= 200 && r.StatusCode <= 299) {
		log.WithFields(log.Fields{
			"status_code": r.StatusCode,
			"post_duration": elapsed.String(),
			"batch_size": len(payloads),
		}).Error("Post api endpoint returned non success status code. Have to retry")
		return errors.New("post api endpoint returned non success status code Have to retry")
	}
	log.WithFields(log.Fields{
		"status_code": r.StatusCode,
		"post_duration": elapsed.String(),
		"batch_size": len(payloads),
	}).Info("Batch post succeeded")
	r.Body.Close()
	return nil
	//return errors.New("testing")
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
