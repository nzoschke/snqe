# snqe

Store, Notify, Queue, Email with S3, SNS, SQS, and SES

GET pre-signed S3 URL
POST a file
Triggers SNS notification and SQS message
GET SQS
POST an SES email

```
$ convox services create ses
$ convox services create sns
$ convox services create sqs --topic $SNS
$ convox services create s3  --topic $SNS

$ convox services link ses
$ convox services link sns
$ convox services link sqs
$ convox services link s3
```