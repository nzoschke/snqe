package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type MessageBody struct {
	Records []Record
}

type Record struct {
	EventSource string `json:"eventSource"`
	EventName   string `json:"eventName"`
	S3          S3     `json:"s3"`
}

type S3 struct {
	Object S3Object `json:"object"`
}

type S3Object struct {
	Key string `json:"key"`
}

func LongPollSQS() {
	svc := sqs.New(session.New(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}))

	for {
		fmt.Printf("sqs.ReceiveMessage WaitTimeSeconds=20\n")
		m, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
			QueueUrl:        aws.String(os.Getenv("SQS_URL")),
			WaitTimeSeconds: aws.Int64(20),
		})

		if err != nil {
			fmt.Printf("sqs.ReceiveMessage error=%q\n", err)
			continue
		}

		if len(m.Messages) == 0 {
			continue
		}

		entries := []*sqs.DeleteMessageBatchRequestEntry{}

		for _, m := range m.Messages {
			e := &sqs.DeleteMessageBatchRequestEntry{
				Id:            m.MessageId,
				ReceiptHandle: m.ReceiptHandle,
			}
			entries = append(entries, e)

			fmt.Printf("sqs.Message.Body=%q\n", *m.Body)

			mb := MessageBody{}
			err := json.Unmarshal([]byte(*m.Body), &mb)

			if err != nil {
				fmt.Printf("json.Unmarshal error=%q\n", err)
				continue
			}

			if len(mb.Records) == 0 {
				continue
			}

			for _, r := range mb.Records {
				fmt.Printf("EventName=%s Key=%s\n", r.EventName, r.S3.Object.Key)
			}
		}

		fmt.Printf("sqs.DeleteMessageBatch NumEntries=%d\n", len(entries))
		_, err = svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
			Entries:  entries,
			QueueUrl: aws.String(os.Getenv("SQS_URL")),
		})

		if err != nil {
			fmt.Printf("sqs.DeleteMessageBatch error=%q\n", err)
			continue
		}
	}
}

func PresignURL(w http.ResponseWriter, r *http.Request) {
	svc := s3.New(session.New(&aws.Config{Region: aws.String(os.Getenv("AWS_REGION"))}))

	u, _ := url.Parse(os.Getenv("S3_URL"))
	bucket := u.Host

	key := fmt.Sprintf("%d", time.Now().UnixNano())
	fmt.Printf("s3.PutObjectRequest.Presign bucket=%s key=%s\n", bucket, key)

	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	str, err := req.Presign(15 * time.Minute)

	if err != nil {
		fmt.Printf("s3.PutObjectRequest.Presign error=%s\n", err)
	}

	fmt.Fprint(w, str)
	return
}

func main() {
	go LongPollSQS()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(PresignURL))
	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", mux)
}
