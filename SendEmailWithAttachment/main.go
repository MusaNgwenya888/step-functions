package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/mail"
	"strings"
	"text/template"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/ses"
)

var (
	htmlTemp string
	tem      *template.Template
)

type OrderInfo struct {
	SupplierID           string `xml:"SupplierID" json:"SupplierID,omitempty"`
	SortKey              string `xml:"SortKey" json:"SortKey,omitempty"`
	OrderID              string `xml:"OrderID" json:"OrderID"`
	BranchID             string `xml:"BranchID" json:"BranchID"`
	AccountID            string `xml:"AccountID" json:"AccountID"`
	AccountName          string `xml:"AccountName" json:"AccountName,omitempty"`
	UserID               string `xml:"UserID" json:"UserID"`
	RepID                string `xml:"RepID" json:"RepID,omitempty"`
	Type                 string `xml:"Type" json:"Type"`
	CreateDate           string `xml:"CreateDate" json:"CreateDate"`
	RequiredByDate       string `xml:"RequiredByDate" json:"RequiredByDate"`
	Reference            string `xml:"Reference" json:"Reference"`
	Comments             string `xml:"Comments" json:"Comments"`
	Route                string `xml:"Route" json:"Route,omitempty"`
	Status               string `xml:"Status" json:"Status,omitempty"`
	Longitude            string `xml:"Longitude" json:"Longitude,omitempty"`
	Latitude             string `xml:"Latitude" json:"Latitude,omitempty"`
	TotalExcl            string `xml:"TotalExcl" json:"TotalExcl,omitempty"`
	RepChangedPrice      string `xml:"RepChangedPrice" json:"RepChangedPrice,omitempty"`
	ClientOrderID        string `xml:"ClientOrderID" json:"ClientOrderID,omitempty"`
	DeliveryName         string `xml:"DeliveryName" json:"DeliveryName"`
	DeliveryAddress1     string `xml:"DeliveryAddress1" json:"DeliveryAddress1,omitempty"`
	DeliveryAddress2     string `xml:"DeliveryAddress2" json:"DeliveryAddress2,omitempty"`
	DeliveryAddress3     string `xml:"DeliveryAddress3" json:"DeliveryAddress3,omitempty"`
	DeliveryPostCode     string `xml:"DeliveryPostCode" json:"DeliveryPostCode,omitempty"`
	DeliveryMethod       string `xml:"DeliveryMethod" json:"DeliveryMethod,omitempty"`
	RouteID              string `xml:"RouteID" json:"RouteID,omitempty"`
	ShipmentID           string `xml:"ShipmentID" json:"ShipmentID,omitempty"`
	PostedToERP          string `xml:"PostedToERP" json:"PostedToERP"`
	ERPOrderNumber       string `xml:"ERPOrderNumber" json:"ERPOrderNumber"`
	ERPStatus            string `xml:"ERPStatus" json:"ERPStatus,omitempty"`
	Email                string `xml:"Email" json:"Email,omitempty"`
	Value                string `xml:"Value" json:"Value,omitempty"`
	Locked               string `xml:"Locked" json:"Locked,omitempty"`
	LockedBy             string `xml:"LockedBy" json:"LockedBy,omitempty"`
	LockedDate           string `xml:"LockedDate" json:"LockedDate,omitempty"`
	PaymentDate          string `xml:"PaymentDate" json:"PaymentDate,omitempty"`
	HaveNotifiedCreator  string `xml:"HaveNotifiedCreator" json:"HaveNotifiedCreator,omitempty"`
	HaveNotifiedCustomer string `xml:"HaveNotifiedCustomer" json:"HaveNotifiedCustomer,omitempty"`
	WorkflowAllowed      string `xml:"WorkflowAllowed" json:"WorkflowAllowed,omitempty"`
	UpdateStockAllowed   string `xml:"UpdateStockAllowed" json:"UpdateStockAllowed,omitempty"`
	UserField01          string `xml:"UserField01" json:"UserField01,omitempty"`
	UserField02          string `xml:"UserField02" json:"UserField02,omitempty"`
	UserField03          string `xml:"UserField03" json:"UserField03,omitempty"`
	UserField04          string `xml:"UserField04" json:"UserField04,omitempty"`
	UserField05          string `xml:"UserField05" json:"UserField05,omitempty"`
	UserField06          string `xml:"UserField06" json:"UserField06,omitempty"`
	UserField07          string `xml:"UserField07" json:"UserField07,omitempty"`
	UserField08          string `xml:"UserField08" json:"UserField08,omitempty"`
	UserField09          string `xml:"UserField09" json:"UserField09,omitempty"`
	UserField10          string `xml:"UserField10" json:"UserField10,omitempty"`
	UserAmount01         string `xml:"UserAmount01" json:"UserAmount01,omitempty"`
	UserAmount02         string `xml:"UserAmount02" json:"UserAmount02,omitempty"`
	UserAmount03         string `xml:"UserAmount03" json:"UserAmount03,omitempty"`
	UserAmount04         string `xml:"UserAmount04" json:"UserAmount04,omitempty"`
	UserAmount05         string `xml:"UserAmount05" json:"UserAmount05,omitempty"`
	UserAmount06         string `xml:"UserAmount06" json:"UserAmount06,omitempty"`
	UserAmount07         string `xml:"UserAmount07" json:"UserAmount07,omitempty"`
	UserAmount08         string `xml:"UserAmount08" json:"UserAmount08,omitempty"`
	UserAmount09         string `xml:"UserAmount09" json:"UserAmount09,omitempty"`
	UserAmount10         string `xml:"UserAmount10" json:"UserAmount10,omitempty"`
	UserAmount11         string `xml:"UserAmount11" json:"UserAmount11,omitempty"`
	UserAmount12         string `xml:"UserAmount12" json:"UserAmount12,omitempty"`
	UserAmount13         string `xml:"UserAmount13" json:"UserAmount13,omitempty"`
	UserAmount14         string `xml:"UserAmount14" json:"UserAmount14,omitempty"`
	UserAmount15         string `xml:"UserAmount15" json:"UserAmount15,omitempty"`
	PaymentTerms         string `xml:"PaymentTerms" json:"PaymentTerms,omitempty"`
	ShipTo               string `xml:"ShipTo" json:"ShipTo,omitempty"`
	Payor                string `xml:"Payor" json:"Payor,omitempty"`
	PayorOrderNumb       string `xml:"PayorOrderNumb" json:"PayorOrderNumb,omitempty"`
	BillTo               string `xml:"BillTo" json:"BillTo,omitempty"`
	Comments2            string `xml:"Comments2" json:"Comments2,omitempty"`
	SMSText              string `xml:"SMSText" json:"SMSText,omitempty"`
	SMSAuth              string `xml:"SMSAuth" json:"SMSAuth,omitempty"`
	SMSNum               string `xml:"SMSNum" json:"SMSNum,omitempty"`
	InvoiceID            string `xml:"InvoiceID" json:"InvoiceID,omitempty"`
	OrderItems           string `xml:"OrderItems>OrderItemInfo" json:"OrderItems"`
	CrossReference       string `xml:"" json:",omitempty"`
	// Approvals            []ApprovalInfo
	Deleted string `xml:"Deleted" json:"Deleted"`
}

func main() {

	// create a new AWS SES client
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	svc := ses.New(sess)

	tempDir := "STEPFUNCTIONS/ORDERSTEMPLATES/"
	tempSupplier := strings.ToUpper("MAHLETEST" + "/")
	custEmail := strings.ToUpper("ORDER")

	download := s3manager.NewDownloader(sess)
	file := &aws.WriteAtBuffer{}
	_, err := download.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String("rapidtradetemplates"),
			Key:    aws.String(tempDir + tempSupplier + custEmail + ".html"),
		})
	if err != nil {
		fmt.Printf("unable to download the s3 file %s\n ", err)
	}
	htmlTemp = string(file.Bytes())

	// tem will hold the email template
	tem, err = template.New("").Parse(htmlTemp)
	if err != nil {
		fmt.Printf("error parsing the html template : %v", err)
	}

	// orderData holds the  data that will be merged in the template
	orderData := OrderInfo{
		SupplierID:           "SupplierID",
		SortKey:              "SortKey",
		OrderID:              "OrderID",
		BranchID:             "BranchID",
		AccountID:            "AccountID",
		AccountName:          "AccountName",
		UserID:               "UserID",
		RepID:                "RepID",
		Type:                 "Type",
		CreateDate:           "formattedDatetime",
		RequiredByDate:       "RequiredByDate",
		Reference:            "Reference",
		Comments:             "Comments",
		Route:                "Route",
		Status:               "Status",
		Longitude:            "Longitude",
		Latitude:             "Latitude",
		TotalExcl:            "excluVat",
		RepChangedPrice:      "RepChangedPrice",
		ClientOrderID:        "ClientOrderID",
		DeliveryName:         "DeliveryName",
		DeliveryAddress1:     "DeliveryAddress1",
		DeliveryAddress2:     "DeliveryAddress2",
		DeliveryAddress3:     "DeliveryAddress3",
		DeliveryPostCode:     "DeliveryPostCode",
		DeliveryMethod:       "DeliveryMethod",
		RouteID:              "RouteID",
		ShipmentID:           "ShipmentID",
		PostedToERP:          "PostedToERP",
		ERPOrderNumber:       "ERPOrderNumber",
		ERPStatus:            "ERPStatus",
		Email:                "Email",
		Value:                "Value",
		Locked:               "Locked",
		LockedBy:             "LockedBy",
		LockedDate:           "LockedDate",
		PaymentDate:          "PaymentDate",
		HaveNotifiedCreator:  "HaveNotifiedCreator",
		HaveNotifiedCustomer: "HaveNotifiedCustomer",
		WorkflowAllowed:      "WorkflowAllowed",
		UpdateStockAllowed:   "UpdateStockAllowed",
		UserField01:          "UserField01",
		UserField02:          "UserField02",
		UserField03:          "UserField03",
		UserField04:          "UserField04",
		UserField05:          "UserField05",
		UserField06:          "UserField06",
		UserField07:          "UserField07",
		UserField08:          "UserField08",
		UserField09:          "UserField09",
		UserField10:          "UserField10",
		UserAmount01:         "TotalWithVat",
		UserAmount02:         "UserAmount02",
		UserAmount03:         "UserAmount03",
		UserAmount04:         "UserAmount04",
		UserAmount05:         "UserAmount05",
		UserAmount06:         "UserAmount06",
		UserAmount07:         "UserAmount07",
		UserAmount08:         "UserAmount08",
		UserAmount09:         "UserAmount09",
		UserAmount10:         "UserAmount10",
		UserAmount11:         "UserAmount11",
		UserAmount12:         "UserAmount12",
		UserAmount13:         "UserAmount13",
		UserAmount14:         "UserAmount14",
		UserAmount15:         "UserAmount15",
		PaymentTerms:         "PaymentTerms",
		ShipTo:               "ShipTo",
		Payor:                "Payor",
		PayorOrderNumb:       "PayorOrderNumb",
		OrderItems:           "OrderItems",
		CrossReference:       "CrossReference",
		Deleted:              "Deleted",
	}

	buff := new(bytes.Buffer)
	tem.Execute(buff, orderData)
	if err != nil {
		fmt.Printf("failed to execute the email : %s", err)
	}

	// define email message
	subject := "Test email with attachment"
	from := mail.Address{Name: "no-reply", Address: "no-reply@rapidtrade.biz"}
	to := mail.Address{Name: "Recipient Name", Address: "musawenkosi@rapidtrade.biz"}
	//body := buff.String()
	attachmentFile := "/Users/musa/Desktop/Work/scripts/Step Functions/SendEmailWithAttachment/attachment.pdf"

	// read attachment file and encode as base64
	attachmentData, err := ioutil.ReadFile(attachmentFile)
	if err != nil {
		panic(err)
	}
	encodedData := base64.StdEncoding.EncodeToString(attachmentData)

	// Create the MIME part for the HTML content
	htmlPart := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n")
	htmlPart += fmt.Sprintf("Content-Transfer-Encoding: base64\r\n\r\n%s", base64.StdEncoding.EncodeToString([]byte(buff.Bytes())))

	// create raw email message
	buffer := bytes.Buffer{}
	buffer.WriteString("Content-Type: multipart/mixed; boundary=example_boundary\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("From: " + from.String() + "\r\n")
	buffer.WriteString("To: " + to.String() + "\r\n")
	buffer.WriteString("Subject: " + subject + "\r\n\r\n")

	buffer.WriteString("--example_boundary\r\n")
	buffer.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	buffer.WriteString(buff.String())

	buffer.WriteString("--example_boundary\r\n")
	buffer.WriteString("Content-Type: application/pdf; name=\"attachment.pdf\"\r\n")
	buffer.WriteString("MIME-Version: 1.0\r\n")
	buffer.WriteString("Content-Disposition: attachment; filename=\"attachment.pdf\"\r\n")
	buffer.WriteString("Content-Transfer-Encoding: base64\r\n\r\n")
	buffer.WriteString(encodedData + "\r\n")

	buffer.WriteString("--example_boundary--\r\n")

	// Combine all parts into a single message
	message := []byte(buffer.String() + htmlPart)

	// send raw email
	input := &ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: message,
		},
	}
	_, err = svc.SendRawEmail(input)
	if err != nil {
		panic(err)
	}
}
