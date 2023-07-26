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

func readFileFromS3() []byte {
	bucket := "mukulmantosh"
	item := "names.json"

	// Create an AWS session
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1")},
	)

	// AWS S3 downloader
	downloader := s3manager.NewDownloader(sess)

	// 4) Download the item from the bucket.
	file, err := os.Create(item)
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
	os.Setenv("SLACK_TOKEN", "xoxb-5633598707492-5631605961058-kk4hnTrsrwqb3deYyJ6dYgEz")
	os.Setenv("CHANNEL_ID", "C05JG6MFGJH")
	api := slack.New(os.Getenv("SLACK_TOKEN"))
	channel := os.Getenv("CHANNEL_ID")

	attachment := slack.Attachment{
		Pretext: "Birthday Wishes 🎉🎉",
		Text:    chatGpt(name),
		Color:   "4287f5",
	}

	_, timestamp, err := api.PostMessage(
		channel,

		slack.MsgOptionAttachments(attachment),
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf("Message sent : %s", timestamp)
}

func chatGpt(name string) string {
	client := openai.NewClient("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")
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

	var birthday []Birthday

	S3File := readFileFromS3()
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
