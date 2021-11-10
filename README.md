# Design
- Create a channel with type LogPayLoad struct
- Create a go-routine called BatchRoutine which accepts this channel as the argument
- Pass the above channel as argument to handler function for /log
- Inside the handler function, de-serialize the input and pass the payload to channel
- BatchRoutine waits for message in channel or timeout inside a for-select-case loop
- When batch size exceeds or timeout occurs, it calls batchProcess function and returns and garbage collects
- processBatch immediately starts a new BatchRoutine as go-routine
- processBatch calls the sendBatch function to post the batch using retry mechanism

## Environment Variables

```shell
ENV BATCH_SIZE=5
ENV BATCH_INTERVAL=1m
ENV POST_ENDPOINT=http://127.0.0.1:8081/test
```

## Usage

Pass the post endpoint as mandatory environment variable. batch size and endpoint has defaults

**Example**
docker build -t webhooktest .
docker run  -p 8080:8080 -e POST_ENDPOINT=http://requestbin.net/r/gbpb57il webhooktest

