# snqe

Store, Notify, Queue, Email with S3, SNS, SQS, and SES

Create backing services. SQS and S3 are connected through an SNS topic.

```
$ convox services create ses
$ convox services create sns
$ convox services create sqs --topic $SNS
$ convox services create s3  --topic $SNS
```

Link services to the app to set S3_URL, etc.

```
$ convox apps create
Creating snqe...

$ convox deploy
Deploying...

$ convox services link ses
$ convox services link sns
$ convox services link sqs
$ convox services link s3
```

The app uses S3_URL to pre-sign an S3 PUT URL to give to a client to upload content.

```
$ APP_URL=$(convox api get /apps/snqe/formation | jq '.[] | select(.name == "web") | .balancer')

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