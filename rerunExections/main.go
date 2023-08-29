package main

import (
	"fmt"
	"log"
	"time"

	utility "steps/Utilities"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sfn"
)

func main() {
	a := utility.NewAwssession("")

	targetDate := time.Date(2023, time.June, 12, 0, 0, 0, 0, time.UTC)

	Live := "arn:aws:states:us-east-1:389633136494:stateMachine:DiscountValues-X7avbHjUiVjB"
	// Test := "arn:aws:states:us-east-1:389633136494:stateMachine:RAPISAMTESTStateMachine-iHrvhFFyHHfR"

	// Specify the ARN of the state machine to re-run failed executions for
	stateMachineArn := Live

	var statusInput string
	fmt.Println("Please enter what type of execution status you want to see:")
	fmt.Println("FAILED")
	time.Sleep(1 * time.Second)
	fmt.Println("SUCCEEDED")
	time.Sleep(1 * time.Second)
	fmt.Println("ABORTED")
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

		descResp, err := a.Sfc.DescribeExecution(descInput)
		if err != nil {
			fmt.Println("Error describing execution:", err)
			continue
		}

		input := *descResp.Input

		fmt.Printf("Here are all executions with the status of '%s' : %s\n", status, input)

		startDate, err := time.Parse(time.RFC3339, execution.StartDate.String())
		if err != nil {
			log.Println("Failed to parse execution start date:", err)
			continue
		}

		if startDate.UTC().Truncate(24 * time.Hour).Equal(targetDate) {
			fmt.Println("Rerunning execution:", *execution.ExecutionArn)

			// Create a new StartExecutionInput with the state machine ARN and input
			startInput := &sfn.StartExecutionInput{
				StateMachineArn: aws.String(stateMachineArn),
				Input:           aws.String(input),
				Name:            aws.String("Retry-" + *execution.Name),
			}

			// Call the StartExecution method to re-run the execution with the specified input
			startResp, err := a.Sfc.StartExecution(startInput)
			if err != nil {
				fmt.Println("Error re-running execution:", err)
				continue
			}

			// Print out the ARN of the new execution that was started
			fmt.Println("New execution started with ARN:", *startResp.ExecutionArn)
		}
	}

	fmt.Println(">>>>>>>>>>>> DONE :-) <<<<<<<<<<<<")
}
