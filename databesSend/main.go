package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/sfn"
	_ "github.com/denisenkom/go-mssqldb"
)

// Contact represents contact information.
type Contact struct {
	ContactID      string        `json:"ContactID"`
	SupplierID     string        `json:"SupplierID"`
	AccountID      string        `json:"AccountID"`
	Contacts       []AccountInfo `json:"Contacts"`
	RepID          string        `json:"RepID"`
	RepIDContactID string        `json:"RepIDContactID"`
	Deleted        bool          `json:"Deleted"`
	Latitude       string        `json:"Latitude"`
	Longitude      string        `json:"Longitude"`
	Counter        string        `json:"Counter"`
}

// AccountInfo represents account information.
type AccountInfo struct {
	SupplierID  string // Change the data type to int
	Counter     string `json:"Counter"`
	RepID       string
	AccountID   string      `xml:"AccountID" json:"AccountID"`
	ContactID   string      `json:"ContactID"`
	Name        string      `json:"Name"`
	Position    string      `json:"Position"`
	TEL         string      `json:"TEL"`
	Mobile      string      `json:"Mobile"`
	Email       string      `json:"Email"`
	UserField1  string      `json:"UserField1"`
	UserField2  string      `json:"UserField2"`
	UserField3  string      `json:"UserField3"`
	UserField4  string      `json:"UserField4"`
	UserField5  string      `json:"UserField5"`
	UserFields  []UserField `json:"UserFields"`
	Deleted     bool        `json:"Deleted"`
	PostedToErp bool        `json:"PostedToErp"`
	ValidForm   bool        `json:"ValidForm"`
}

// UserField represents a user field.
type UserField struct {
	Name  string `json:"Name"`
	Value string `json:"Value"`
}

// Declare contactObjects at the package level
var contactObjects = make(map[string]Contact)

func main() {
	connString := "server=sqli.rapidtrade.biz;user id=beston152;password=pass@word3;database=rapidtrade"
	db, err := sql.Open("sqlserver", connString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected!")

	query := "SELECT * FROM Contacts WHERE supplierid = 'DWS'"

	rows, err := db.Query(query)
	if err != nil {
		log.Fatal("failed to run query", err)
	}
	defer rows.Close()

	// Create a slice to store retrieved items
	var contactList []AccountInfo

	for rows.Next() {
		var c AccountInfo
		if err := rows.Scan(
			&c.SupplierID,
			&c.AccountID,
			&c.Counter,
			&c.ContactID,
			&c.Name,
			&c.Position,
			&c.TEL,
			&c.Mobile,
			&c.Email,
			&c.UserField1,
			&c.UserField2,
			&c.UserField3,
			&c.UserField4,
			&c.PostedToErp,
			&c.UserField5,
			&c.RepID,
		); err != nil {
			log.Fatal("failed to scan successfully", err)
		}

		c.ContactID = genGuid()
		initializeContactInfo(&c)
		initializeAccountInfo(&c)

		// Append the retrieved item to the contactList slice
		contactList = append(contactList, c)
	}

	// Marshal the contactList to JSON
	jsonData, err := json.Marshal(contactList)
	if err != nil {
		log.Fatal("failed to marshal data to JSON", err)
	}

	// log.Fatal("this is json data", contactList)

	// Create an AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-1"), // Specify your AWS region
	})

	if err != nil {
		log.Fatal(err)
	}

	// Create an S3 service client
	s3Svc := s3.New(sess)

	// Specify the S3 bucket and object key
	bucketName := "rapidtradeinbox"
	objectKey := "BulkDWSContactFile.json" // Modify with your desired S3 object key

	// Upload the JSON data to S3
	_, err = s3Svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(jsonData),
	})

	if err != nil {
		log.Fatal("failed to upload JSON data to S3", err)
	}

	println("JSON data uploaded to S3 successfully")

	// Specify the name of your Step Functions state machine
	stateMachineName := "arn:aws:states:us-east-1:389633136494:stateMachine:ContactsTest-kxuFn8Y1PEms"

	// Specify the input for the Step Functions execution
	input := "{\"s3Bucket\": \"" + bucketName + "\", \"S3File\": \"" + objectKey + "\"}"

	// Create an SFN (Step Functions) service client
	sfnSvc := sfn.New(sess)

	// Start the Step Functions execution
	_, err = sfnSvc.StartExecution(&sfn.StartExecutionInput{
		StateMachineArn: aws.String(stateMachineName),             // Modify with your state machine ARN
		Name:            aws.String("DWS_CONTACTS_BULK_RUN_LIVE"), // Modify with a unique execution name
		Input:           aws.String(input),
	})

	if err != nil {
		log.Fatal("failed to start Step Functions execution", err)
	}

	println("Step Functions execution started successfully")
}

func initializeContactInfo(obj *AccountInfo) {
	if _, exists := contactObjects[obj.AccountID]; exists {
		return
	}
	cIobj := Contact{
		SupplierID:     "DWS",
		AccountID:      obj.AccountID,
		RepID:          "",
		RepIDContactID: obj.RepID + "|" + obj.ContactID,
		Deleted:        obj.Deleted,
		Latitude:       "",
		Longitude:      "",
		Counter:        obj.Counter,
	}
	contactObjects[cIobj.AccountID] = cIobj
}

func initializeAccountInfo(obj *AccountInfo) {
	// Initialize the UserFields array
	obj.UserFields = []UserField{
		{Name: "UserField1", Value: ""},
		{Name: "UserField2", Value: ""},
		{Name: "UserField3", Value: ""},
		{Name: "UserField4", Value: ""},
		{Name: "UserField5", Value: ""},
	}
}

// func post(c []Contact) error {
// 	url := "http://rapi.rapidtradews.com/post/contacts"
// 	username := "DWSTEST"
// 	password := "PASSWORD"

// 	auth := username + ":" + password
// 	authEncoded := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

// 	jsonData, err := json.Marshal(c)
// 	if err != nil {
// 		return err
// 	}

// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
// 	if err != nil {
// 		return err
// 	}

// 	req.Header.Set("Authorization", authEncoded)
// 	req.Header.Set("Content-Type", "application/json")
// 	req.Header.Set("Max-Content-Length", "524288900")

// 	client := &http.Client{}
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
// 		fmt.Println("Data posted successfully.")
// 	} else {
// 		fmt.Println("Error posting data. Status code:", resp.StatusCode)
// 	}

// 	return nil
// }

func genGuid() string {
	date := time.Now()
	yy := date.Format("06")
	JJJ := getJulianDay()
	dd := fmt.Sprintf("%02d", date.Day())
	hh := fmt.Sprintf("%02d", date.Hour())
	mm := fmt.Sprintf("%02d", date.Minute())
	ss := fmt.Sprintf("%02d", date.Second())
	rr := fmt.Sprintf("%02d", rand.Intn(100))
	return yy + JJJ + dd + hh + mm + ss + rr
}

func getJulianDay() string {
	date := time.Now()
	// onejan := time.Date(date.Year(), time.January, 1, 0, 0, 0, 0, time.UTC)
	JJJ := (date.YearDay() + 1)
	return fmt.Sprintf("%03d", JJJ)
}
