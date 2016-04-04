# snqe

Store, Notify, Queue, Email with S3, SNS, SQS, and SES

This demonstrates how Convox simplifies managing AWS services:

* Provision an SQS queue with a single command
* Provision an S3 bucket configured to send notifications to the SQS queue with a single command
* Get and set S3_URL and SQS_URL in an app environment
* Interact with S3 and SQS in the app

This is a work in progress. SNS, SES and service/app linking are not yet implemented.

## Setup

Create backing SQS and S3 services

```
$ convox services create sqs --name mysqs
$ convox services create s3  --name mys3 --queue mysqs
```

Create an app and set service URLs

```
$ convox apps create
Creating app snqe... CREATING

$ convox env set SQS_URL=$(convox api get /services/mysqs | jq -r '.exports.URL')
$ convox env set S3_URL=$(convox api get /services/mys3 | jq -r '.exports.URL')

$ convox deploy
Deploying...
```

## App Logic

The app uses S3_URL to pre-sign an S3 PUT URL to give to a client to upload content.

```
$ APP_URL=$(convox api get /apps/snqe/formation | jq -r '.[] | select(.name == "web") | .balancer')

$ curl -sS -i -X PUT -T Readme.md $(curl -s $APP_URL)
HTTP/1.1 100 Continue

HTTP/1.1 200 OK
x-amz-id-2: +QH4nodB1aJc7UOZ9j+R7aWs0HonDY0f8QO29WQoryMtWtbgQguN82G6kufJHr5q/Dy5QDmjnyM=
x-amz-request-id: 5E06AA82313B40E8
Date: Sat, 19 Mar 2016 23:38:08 GMT
ETag: "8d2fc2b4ec70f822d49680232ab420bd"
Content-Length: 0
Server: AmazonS3
```

The app gets and acts on S3 events through SQS messages

```
$ convox logs
sqs.ReceiveMessage WaitTimeSeconds=20
s3.PutObjectRequest.Presign bucket=convox-s3-1621 key=1458430165990080268
EventName=ObjectCreated:Put Key=1458430165990080268
sqs.DeleteMessageBatch NumEntries=1
sqs.ReceiveMessage WaitTimeSeconds=20
```