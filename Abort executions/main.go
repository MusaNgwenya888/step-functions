package main

import (
	"fmt"
	"time"

	utility "StepFunctions/Utilities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func main() {
	a := utility.NewAwssession("")
	fmt.Println(">>>>>>>>>>>> Starting :-) <<<<<<<<<<<<")
	// ActivityLive := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIES-LI0zVfbremKg"
	// ActivityTest := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIESTEST-gqNiSvxWmzQL"
	// OrdersLive := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMStateMachine-1YwcU3XmIcZ5"
	OrdersTest := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMTESTStateMachine-iHrvhFFyHHfR"

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := OrdersTest

	var statusInput string
	fmt.Println("Please enter what type of execution status you want to see:")
	fmt.Println("FAILED")
	time.Sleep(1 * time.Second)
	fmt.Println("SUCCEEDED")
	time.Sleep(1 * time.Second)
	fmt.Println("RUNNING")
	fmt.Println("Type now...")
	fmt.Scanln(&statusInput)
	status := statusInput

	// List the failed executions for the state machine
	listInput := &sfn.ListExecutionsInput{
		StateMachineArn: aws.String(stateMachineArn),
		StatusFilter:    aws.String(status),
	}

	resp, err := a.Sfc.ListExecutions(listInput)
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
		_, err := a.Sfc.StopExecution(stopInput)
		if err != nil {
			fmt.Println("Error stopping execution: ", err)
		} else {
			fmt.Println("Stopped execution: ", *exec.ExecutionArn)
		}
	}

	fmt.Println(">>>>>>>>>>>> DONE :-) <<<<<<<<<<<<")
}
