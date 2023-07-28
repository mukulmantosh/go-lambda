# SlackBot with Go & OpenAI

![background](./misc/background.png)


## Dependencies

Install the below packages

```bash
$ go get github.com/sashabaranov/go-openai
$ go get github.com/aws/aws-sdk-go
$ go get -u github.com/slack-go/slack
```

### Building Docker Image
```
$ docker build -t go-lambda:latest .
```


### Deploying Image

Follow the instructions provided over here https://docs.aws.amazon.com/lambda/latest/dg/go-image.html
### Creating Function

```bash
aws lambda create-function --function-name <FUNCTION-NAME> --package-type Image --code <ECR-IMAGE-URL> --role <AWS-ROLE-NAME>
```


### Updating Function

```bash
aws lambda update-function-code --function-name <FUNCTION-NAME> --image-uri <ECR-IMAGE-URL>
```

