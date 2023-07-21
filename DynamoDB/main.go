package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/s3"
)

func main() {
	// Create a new AWS session using your credentials and desired region.
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Change this to your desired AWS region
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	// Download the JSON file from S3.
	bucket := "rapidtradeinbox"
	objectKey := "products.json"

	s3Svc := s3.New(sess)
	result, err := s3Svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Fatalf("Failed to download JSON file from S3: %v", err)
	}
	defer result.Body.Close()

	// Read and parse the JSON data.
	jsonData, err := ioutil.ReadAll(result.Body)
	if err != nil {
		log.Fatalf("Failed to read JSON data: %v", err)
	}

	// Define a struct that matches the structure of the JSON data.
	// This is an example, adjust it according to your JSON data structure.
	type Product struct {
		ProductID string `json:"ProductID"`
	}

	// Unmarshal the JSON data into a slice of Product structs.
	var products []Product
	err = json.Unmarshal(jsonData, &products)
	if err != nil {
		log.Fatalf("Failed to unmarshal JSON data: %v", err)
	}

	// Create a new DynamoDB client.
	dynamoDBSvc := dynamodb.New(sess)

	// Set the table name for the items you want to delete.
	tableName := "TradeProducts"

	// Perform the delete operations using the product IDs from the JSON file.
	for _, product := range products {
		params := &dynamodb.DeleteItemInput{
			TableName: aws.String(tableName),
			Key: map[string]*dynamodb.AttributeValue{
				"SupplierID": {
					S: aws.String("RAPISAM_MD_BIG"), // Replace with the actual supplier ID.
				},
				"ProductID": {
					S: aws.String(product.ProductID), // Use the product ID from the JSON data.
				},
			},
		}

		_, err := dynamoDBSvc.DeleteItem(params)
		if err != nil {
			log.Fatalf("Failed to delete item: %v", err)
		}

		fmt.Printf("Item with ProductID '%s' deleted successfully.\n", product.ProductID)
	}

	fmt.Println("All items deleted successfully.")
}
