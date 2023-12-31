package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	openai "github.com/sashabaranov/go-openai"
	"github.com/slack-go/slack"
	"log"
	"os"
	"time"
)

type Birthday struct {
	Name string `json:"name"`
	DOB  string `json:"dob"`
}

func readFileFromS3(bucket string, item string) []byte {
	// Create an AWS session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	// AWS S3 downloader
	downloader := s3manager.NewDownloader(sess)

	// 4) Download the item from the bucket.
	file, err := os.Create("/tmp/" + item)
	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(item),
		})

	if err != nil {
		log.Fatalf("Unable to download item %q, %v", item, err)
	}

	data, err := os.ReadFile(file.Name())
	if err != nil {
		fmt.Println("File reading error", err)
	}
	return []byte(data)
}

func sendSlackMessage(name string) {
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := []string{os.Getenv("CHANNEL_ID")}

	attachment := slack.Attachment{
		Pretext: "Birthday Wishes 🎉🎉",
		Text:    chatGpt(name),
		Color:   "4287f5",
	}

	_, timestamp, err := api.PostMessage(
		channel[0],
		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent : %s", timestamp)

	// Send Birthday Photo
	readFileFromS3(os.Getenv("BUCKET_NAME"), "birthday.jpg")

	params := slack.FileUploadParameters{
		Channels: channel,
		File:     "/tmp/birthday.jpg",
	}

	_, err = api.UploadFile(params)
	if err != nil {
		fmt.Printf("Unable to upload file to Slack %s\n", err)
	}

}

func chatGpt(name string) string {
	client := openai.NewClient(os.Getenv("OPENAI_TOKEN"))
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Provide me a birthday message for person name " + name + "under 15 words and not to be repeated",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return ""
	}

	return resp.Choices[0].Message.Content
}

func handler(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	currentTime := time.Now()
	currentDate := currentTime.Format("2006-01-02")

	bucket := os.Getenv("BUCKET_NAME")
	item := os.Getenv("FILE_NAME")

	var birthday []Birthday

	S3File := readFileFromS3(bucket, item)
	json.Unmarshal(S3File, &birthday)

	for _, v := range birthday {
		if v.DOB == currentDate {
			sendSlackMessage(v.Name)
		}
	}
	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "\"Message Sent!\"",
	}
	return response, nil
}

func main() {
	lambda.Start(handler)
}
