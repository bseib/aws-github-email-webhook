package main

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	//go get -u github.com/aws/aws-lambda-go
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	//go get -u github.com/aws/aws-sdk-go
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// The character encoding for the email.
const CharSet = "UTF-8"

func getAwsRegion() string {
	return os.Getenv("SES_AWS_REGION")
}

func getGithubWebhookSecret() string {
	return os.Getenv("GITHUB_WEBHOOK_SECRET")
}

func getSender() string {
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	return os.Getenv("SENDER")
}

func getCommaRecipients() string {
	return os.Getenv("RECIPIENTS")
}
func getRecipients() []*string {
	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	s := strings.Split(getCommaRecipients(), ",")
	result := make([]*string, len(s))
	for i, _ := range s {
		result[i] = &s[i]
	}
	return result
}

func isKnownEmail(email string) bool {
	for _, knownEmail := range getRecipients() {
		if *knownEmail == email {
			return true
		}
	}
	return false
}

// checkMAC reports whether messageMAC is a valid HMAC tag for message.
func checkMAC(message, messageMAC, key []byte) bool {
	mac := hmac.New(sha1.New, key)
	mac.Write(message)
	expectedMAC := mac.Sum(nil)
	return hmac.Equal(messageMAC, expectedMAC)
}

func main() {
	lambda.Start(Handler)
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	var xHubSignature string = request.Headers["X-Hub-Signature"]
	fmt.Println("X-Hub-Signature: " + xHubSignature)
	// signature begins with these 5 characters 'sha1='
	payloadHmac, _ := hex.DecodeString(xHubSignature[5:])
	payloadBody := request.Body
	isValidHmac := checkMAC([]byte(payloadBody), []byte(payloadHmac), []byte(getGithubWebhookSecret()))
	if !isValidHmac {
		fmt.Println("  `--> signature is invalid")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusForbidden,
			Body:       "",
		}, nil
	} else {
		fmt.Println("  `--> signature is valid")
	}

	var githubEventHeader string = request.Headers["X-GitHub-Event"]
	if "ping" == githubEventHeader {
		fmt.Println("responding OK to github ping event")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       "",
		}, nil
	}

	var githubPushEvent PushEvent
	err := json.Unmarshal([]byte(payloadBody), &githubPushEvent)
	if nil != err {
		fmt.Println(err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       err.Error(),
		}, nil
	}
	fmt.Println(githubPushEvent)

	senderEmail := getSender()
	if nil != githubPushEvent.Pusher && isKnownEmail(*githubPushEvent.Pusher.Email) {
		senderEmail = *githubPushEvent.Pusher.Email
		fmt.Printf("Pusher %s is a known email. Will use it for Sender.\n", senderEmail)
	}

	// Create a new session in the us-west-2 region.
	// Replace us-west-2 with the AWS Region you're using for Amazon SES.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(getAwsRegion())},
	)

	// Create an SES session.
	svc := ses.New(sess)

	// Assemble the email.
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{},
			ToAddresses: getRecipients(),
		},
		Message: &ses.Message{
			Body: &ses.Body{
				// Html: &ses.Content{
				// 	Charset: aws.String(CharSet),
				// 	Data:    aws.String(HtmlBody),
				// },
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(BuildEmailMessageBody(githubPushEvent)),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String(BuildEmailSubject(githubPushEvent)),
			},
		},
		Source: aws.String(senderEmail),
		// Uncomment to use a configuration set
		//ConfigurationSetName: aws.String(ConfigurationSet),
	}

	// Attempt to send the email.
	result, err := svc.SendEmail(input)

	// Display error messages if they occur.
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ses.ErrCodeMessageRejected:
				fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
			case ses.ErrCodeMailFromDomainNotVerifiedException:
				fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
			case ses.ErrCodeConfigurationSetDoesNotExistException:
				fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}

		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       err.Error(),
		}, nil
	}

	fmt.Println("Email Sent to: " + getCommaRecipients())
	fmt.Println(result)

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "",
	}, nil
}
