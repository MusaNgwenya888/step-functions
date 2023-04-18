package utility

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sfn"
)

type Awssession struct {
	Config      aws.Config
	Session     *session.Session
	Profile     string
	Region      string
	ErrorLogger *log.Logger
	InfoLogger  *log.Logger
	Svc         *dynamodb.DynamoDB
	S3c         *s3.S3
	Sfc         *sfn.SFN
}

// NewAwssession constructs our aws session
func NewAwssession(region string) *Awssession {
	// checking if the region was passed down , if the region is not passed down it will get from env variable
	if region == "" {
		region = "us-east-1"
	}
	a := Awssession{
		Region:      region,
		ErrorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate+log.Ltime+log.Lshortfile),
		InfoLogger:  log.New(os.Stderr, "INFO: ", log.Ldate+log.Ltime+log.Lshortfile),
	}
	profile := os.Getenv("AWS_PROFILE")
	var err error
	// chooses the session connection based on if there is a value in profile
	if profile == "" {
		a.Session, err = session.NewSession(&aws.Config{
			Region: aws.String(region),
		})
	} else {
		conf := aws.Config{Region: aws.String(region)}
		a.Session, err = session.NewSessionWithOptions(session.Options{
			Config:  conf,
			Profile: profile,
		})
	}
	if err != nil {
		log.Fatal(err)
	}
	// setting the service configuration for service clients
	a.Config = aws.Config{Region: aws.String(region)}
	if err != nil {
		log.Fatal(err)
	}
	// creating a new instance of the dynamoDB client with a session
	a.Svc = dynamodb.New(a.Session, &a.Config)
	//creating a new instance of the s3 client with a session
	a.S3c = s3.New(a.Session, &a.Config)
	// session for the step functions
	a.Sfc = sfn.New(a.Session, &a.Config)

	return &a
}

// LogF logs a message similar to PrintF
func (a *Awssession) LogF(loguid string, message string, args ...interface{}) string {
	message = fmt.Sprintf("%s %s", loguid, message)
	fmessage := fmt.Sprintf(message, args...)
	a.InfoLogger.Printf(fmessage)
	return fmessage
}

// LogJ logs a message similar to PrintF but for json to be displayed
func (a *Awssession) LogJ(id string, obj interface{}) string {
	b, err := json.Marshal(obj)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", id, string(b))
}

// LogE logs an error with an additional message, it will also return message for further processing
func (a *Awssession) LogE(loguid string, message string, err error) error {
	message = fmt.Sprintf("%s %s %v", loguid, message, err)
	a.ErrorLogger.Printf(message)
	return errors.New(message)

}
