# SlackBot with Go & OpenAI

![background](./misc/background.png)


```bash
go get github.com/sashabaranov/go-openai

go get github.com/aws/aws-sdk-go

go get -u github.com/slack-go/slack



aws lambda create-function --function-name hello-world --package-type Image --code ImageUri=254501641575.dkr.ecr.ap-south-1.amazonaws.com/go-lambda:latest --role arn:aws:iam::254501641575:role/LambdaRole

// Update
aws lambda update-function-code --function-name hello-world --image-uri 254501641575.dkr.ecr.ap-south-1.amazonaws.com/go-lambda:latest

```

