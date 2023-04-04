package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func main() {
	// ActivityLive := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIES-LI0zVfbremKg"
	// ActivityTest := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIESTEST-gqNiSvxWmzQL"
	OrdersLive := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMStateMachine-1YwcU3XmIcZ5"
	// OrdersTest := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMTESTStateMachine-iHrvhFFyHHfR"

	// Create a new session with the AWS SDK for Go
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create a new StepFunctions client
	svc := sfn.New(sess)

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := OrdersLive
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

	var typeInput string
	fmt.Printf("Are you sure you want to re run all %s executions ?\n", status)
	fmt.Println("Please enter y or n:")
	fmt.Scanln(&typeInput)
	if typeInput == "y" {
		fmt.Println("STARTING TO LIST ALL EXECUTIONS :-D")
	} else if typeInput == "n" {
		fmt.Println("CANCELLED :-(")
		return
	} else {
		fmt.Println("Invalid input. :-/")
		return
	}

	// Loop through the executions and retrieve the input for each one
	for _, execution := range resp.Executions {
		// Specify the ARN of the execution to retrieve the input for
		executionArn := *execution.ExecutionArn

		// Call the DescribeExecution method to retrieve information about the execution
		descInput := &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		}

		descResp, err := svc.DescribeExecution(descInput)
		if err != nil {
			fmt.Println("Error describing execution:", err)
			continue
		}

		input := *descResp.Input

		fmt.Printf("Here are all executions with the status of '%s' : %s\n", status, input)

		// // Create a new StartExecutionInput with the state machine ARN and input
		// startInput := &sfn.StartExecutionInput{
		// 	StateMachineArn: aws.String(stateMachineArn),
		// 	Input:           aws.String(input),
		// 	Name:            aws.String("Retry-" + *execution.Name),
		// }

		// // Call the StartExecution method to re-run the execution with the specified input
		// startResp, err := svc.StartExecution(startInput)
		// if err != nil {
		// 	fmt.Println("Error re-running execution:", err)
		// 	continue
		// }

		// // Print out the ARN of the new execution that was started
		// fmt.Println("New execution started with ARN:", *startResp.ExecutionArn)
	}

	fmt.Println(">>>>>>>>>>>> DONE :-) <<<<<<<<<<<<")
}
