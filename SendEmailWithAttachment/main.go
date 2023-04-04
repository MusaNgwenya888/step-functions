package main

import (
	"bytes"
	"context"
	"fmt"
	"html/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

func main() {
	// Initialize AWS SES session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2"),
	})
	if err != nil {
		panic(err)
	}

	// Create SES client
	svc := ses.New(sess)

	// Load HTML template
	htmlTemplate, err := template.ParseFiles("email_template.html")
	if err != nil {
		panic(err)
	}

	// Render HTML template with data
	var rendered bytes.Buffer
	data := map[string]interface{}{
		"Name": "John Doe",
	}
	err = htmlTemplate.Execute(&rendered, data)
	if err != nil {
		panic(err)
	}

	// Load attachment file
	// attachmentContent := []byte("This is the attachment content")
	// attachmentFileName := "attachment.txt"

	// // Convert attachment content to base64 encoded string
	// attachmentContentBase64 := base64.StdEncoding.EncodeToString(attachmentContent)

	// Create SES input
	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String("recipient@example.com"),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data: aws.String(rendered.String()),
				},
			},
			Subject: &ses.Content{
				Data: aws.String("Test email with attachment"),
			},
		},
		Source: aws.String("sender@example.com"),
	}

	// Add attachment to SES input
	// attachment := &ses.Attachment{
	// 	Data:         []byte(attachmentContentBase64),
	// 	ContentType:  aws.String("text/plain"),
	// 	Filename:     aws.String(attachmentFileName),
	// 	ContentDispo: aws.String("attachment"),
	// }
	// input.Message.Attachments = []*ses.Attachment{attachment}

	// Send email
	_, err = svc.SendEmailWithContext(context.Background(), input)
	if err != nil {
		panic(err)
	}

	fmt.Println("Email sent with attachment")
}
