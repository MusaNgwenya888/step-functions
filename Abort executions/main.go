package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func main() {
	fmt.Println(">>>>>>>>>>>> Starting :-) <<<<<<<<<<<<")
	// ActivityLive := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIES-LI0zVfbremKg"
	// ActivityTest := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIESTEST-gqNiSvxWmzQL"
	// OrdersLive := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMStateMachine-1YwcU3XmIcZ5"
	OrdersTest := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMTESTStateMachine-iHrvhFFyHHfR"

	// Create a new session with the AWS SDK for Go
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new StepFunctions client
	svc := sfn.New(sess)

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := OrdersTest
	status := "RUNNING"

	// List the failed executions for the state machine
	listInput := &sfn.ListExecutionsInput{
		StateMachineArn: aws.String(stateMachineArn),
		StatusFilter:    aws.String(status),
	}

	resp, err := svc.ListExecutions(listInput)
	if err != nil {
		fmt.Println("Error listing executions:", err)
		return
	}

	// Loop through the list of running executions and call the stop execution API for each one
	for _, exec := range resp.Executions {
		stopInput := &sfn.StopExecutionInput{
			Cause:        aws.String("Aborting execution"),
			Error:        aws.String("Aborted"),
			ExecutionArn: exec.ExecutionArn,
		}
		_, err := svc.StopExecution(stopInput)
		if err != nil {
			fmt.Println("Error stopping execution: ", err)
		} else {
			fmt.Println("Stopped execution: ", *exec.ExecutionArn)
		}
	}

	fmt.Println(">>>>>>>>>>>> DONE :-) <<<<<<<<<<<<")
}
