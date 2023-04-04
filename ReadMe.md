# Step Functions

## Table of Contents 
1. Documentation.
    - Requirements
    - Set Up Process
2. Functions
    - Abort Executions
    - List Executions
    - Re Run Executions

## 1. Documentation

### Requirements
- Install AWS CLI (https://docs.aws.amazon.com/cli/latest/userguide/install-cliv2-windows.html)
- Install Golang [https://golang.org/dl/]

### Set Up Process
```bash
go mod tidy
```
## 2. Functions
### Abort Executions
- This function aborts step functions executions of your choice, choices :
    - FAILED
    - SUCEEDED
    - RUNNING

### List Executions
- This functions lists all executions to your choice, choices :
    - FAILED
    - SUCEEDED
    - RUNNING
    - ABORTED

### Re Run Executions
- This functions Re runs all executions to your choice, choices :
    - FAILED
    - SUCEEDED
    - RUNNING
    - ABORTED