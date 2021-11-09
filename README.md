## Environment Variables

ENV BATCH_SIZE=5
ENV BATCH_INTERVAL=1m
ENV POST_ENDPOINT=http://127.0.0.1:8081/test

## Usage

Pass the post endpoint as mandatory environment variable. batch size and endpoint has defaults

**Example**
docker build -t webhooktest .
docker run  -p 8080:8080 -e POST_ENDPOINT=http://requestbin.net/r/gbpb57il webhooktest

