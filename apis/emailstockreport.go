package apis

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/sendgrid/sendgrid-go"
)

//EmailRequest request body for a dynamic template email
type EmailRequest struct {
	Message    []Personalizations `json:"personalizations"`
	From       map[string]string  `json:"from"`
	Content    []ContentType      `json:"content"`
	TemplateID string             `json:"template_id"`
}

//Personalizations mapped values for sending template data
type Personalizations struct {
	UserInfo []UserEmails `json:"to"`
	Message  Report       `json:"dynamic_template_data"`
}

//UserEmails user information
type UserEmails struct {
	EmailAddress string `json:"email"`
	Name         string `json:"name"`
}

//ContentType headers for request
type ContentType struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

//Report contains stock data information
type Report struct {
	StockInfo []StockData `json:"stock_data"`
}

//GetStockReportFromTopic listens to a topic
//if a message is published to topic, fetches the data and emails a stock report
func GetStockReportFromTopic(ctx context.Context, msg PubSubMessage) error {
	var report StockReport

	err := json.Unmarshal(msg.Data, &report)
	if err != nil {
		log.Fatalf("error unmarshaling payload from topic: %v", err)
		return err
	}

	err = SendDynamicTemplateEmail(report)
	if err != nil {
		return err
	}

	return nil
}

//SendDynamicTemplateEmail ...
func SendDynamicTemplateEmail(report StockReport) error {
	//Get a sendgrid request object
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")

	//Get the marshaled populated request body
	params, _ := json.Marshal(populateRequestBody(report))
	//populate sendgrid request object parameters
	request.Method = "POST"
	request.Body = params

	//send email request
	resp, err := sendgrid.API(request)
	if err != nil {
		return err
	}

	log.Printf("StatusCode: %v\n, ResponseBody: %v\n, ResponseHeader: %v\n", resp.StatusCode, resp.Body, resp.Headers)

	return nil
}

func populateRequestBody(report StockReport) EmailRequest {
	var req EmailRequest

	fullName := report.User.FirstName + " " + report.User.LastName
	req.Content = []ContentType{{Type: "text/html", Value: "1.1"}}
	req.Message = []Personalizations{{[]UserEmails{{EmailAddress: report.User.Email, Name: fullName}}, Report{StockInfo: report.StockReport}}}
	req.From = map[string]string{"email": os.Getenv("FROM_EMAIL")}
	req.TemplateID = os.Getenv("STOCK_REPORT_TEMPLATE_ID")

	return req
}
