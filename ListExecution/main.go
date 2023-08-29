package main

import (
	"fmt"
	"time"

	utility "steps/Utilities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func main() {
	a := utility.NewAwssession("")
	// ActivityLive := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIES-LI0zVfbremKg"
	// ActivityTest := "arn:aws:states:us-east-1:389633136494:stateMachine:ACTIVITIESTEST-gqNiSvxWmzQL"
	OrdersLive := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMStateMachine-1YwcU3XmIcZ5"
	// OrdersTest := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMTESTStateMachine-iHrvhFFyHHfR"

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := OrdersLive

	var statusInput string
	fmt.Println("Please enter what type of execution status you want to see:")
	fmt.Println("FAILED")
	time.Sleep(1 * time.Second)
	fmt.Println("SUCCEEDED")
	time.Sleep(1 * time.Second)
	fmt.Println("ABORTED")
	time.Sleep(1 * time.Second)
	fmt.Println("RUNNING")
	time.Sleep(1 * time.Second)
	fmt.Println("TIMED_OUT")
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

	// Loop through the executions and retrieve the input for each one
	for _, execution := range resp.Executions {
		// Specify the ARN of the execution to retrieve the input for
		executionArn := *execution.ExecutionArn

		// Call the DescribeExecution method to retrieve information about the execution
		descInput := &sfn.DescribeExecutionInput{
			ExecutionArn: aws.String(executionArn),
		}

		descResp, err := a.Sfc.DescribeExecution(descInput)
		if err != nil {
			fmt.Println("Error describing execution:", err)
			continue
		}

		input := *descResp.Input
		fmt.Printf("Status : %s,\n Input : %s\n\n", status, input)
	}
	fmt.Println(">>>>>>>>>>>>>>>>>>>> DONE :-) <<<<<<<<<<<<<<<<<<<<<<<<<<<")
}
