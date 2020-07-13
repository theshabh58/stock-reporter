package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var tpl *template.Template
var emailAuth smtp.Auth

//DymanicRequest ...
type DymanicRequest struct {
	Personalizations []DataValues      `json:"personalizations"`
	From             map[string]string `json:"from"`
	Content          []ContentType     `json:"content"`
	TemplateID       string            `json:"template_id"`
}

type DataValues struct {
	UserEmails []UserData        `json:"to"`
	Message    StockReportValues `json:"dynamic_template_data"`
}

type UserData struct {
	EmailAddress string `json:"email"`
	Name         string `json:"name"`
}
type ContentType struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type StockReportValues struct {
	StockInfo []StockData `json:"stock_data"`
}

func init() {
	tpl = template.Must(template.ParseFiles("email-template.html"))

}

//GetStockReportFromTopic listens to a topic
//if a message is published to topic, fetches the data and emails a stock report
func GetStockReportFromTopic(ctx context.Context, msg []byte) error {
	var report StockReport

	err := json.Unmarshal(msg, &report)
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

//SendEmail Uses SMTP to send an email using MailTrap API
func SendEmail(report StockReport) error {
	//buffer to store the templated stock values
	buffer := new(bytes.Buffer)

	//populate the report to buffer
	err := tpl.Execute(buffer, report)
	if err != nil {
		return err
	}

	//setup fields to send email through sendgrid
	to := mail.NewEmail(report.User.FirstName, report.User.Email)
	subject := "Stock Report"
	from := mail.NewEmail("Stock-Report", os.Getenv("FROM_EMAIL"))
	htmlContent := buffer.String()
	email := mail.NewSingleEmail(from, subject, to, "None NEeded", htmlContent)

	//setup sendgrid client
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	resp, err := client.Send(email)
	if err != nil {
		return err
	}
	log.Printf("StatusCode: %v\n, ResponseBody: %v\n, ResponseHeader: %v\n", resp.StatusCode, resp.Body, resp.Headers)
	return nil
}

//SendDynamicTemplateEmail ...
func SendDynamicTemplateEmail(report StockReport) error {
	request := sendgrid.GetRequest(os.Getenv("SENDGRID_API_KEY"), "/v3/mail/send", "https://api.sendgrid.com")

	params, _ := json.Marshal(getPopulatedRequest(report))
	fmt.Printf("Request Body %s", params)
	request.Method = "POST"
	request.Body = params

	resp, err := sendgrid.API(request)
	if err != nil {
		return err
	}
	log.Printf("StatusCode: %v\n, ResponseBody: %v\n, ResponseHeader: %v\n", resp.StatusCode, resp.Body, resp.Headers)
	return nil
}

func getPopulatedRequest(report StockReport) DymanicRequest {
	var req DymanicRequest
	req.Content = []ContentType{{Type: "text/html", Value: "1.1"}}
	req.Personalizations = []DataValues{{UserEmails: []UserData{{Name: report.User.FirstName, EmailAddress: report.User.Email}}, Message: StockReportValues{report.StockReport}}}
	req.From = map[string]string{"email": os.Getenv("FROM_EMAIL")}
	req.TemplateID = os.Getenv("TEMPLATE_ID")
	return req
}
