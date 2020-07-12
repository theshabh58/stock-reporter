package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

var tpl *template.Template
var emailAuth smtp.Auth

//EmailParameters ...
type EmailParameters struct {
	To      []string
	From    string
	Subject string
	Message []byte
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

	err = SendEmail(report)
	if err != nil {
		return err
	}

	return nil
}

//SendEmail Uses SMTP to send an email using MailTrap API
func SendEmail(report StockReport) error {
	//buffer to store the templated stock values
	buffer := new(bytes.Buffer)
	// headers := "MIME-version: 1.0;\nContext-Type: text/html;"
	// buffer.Write([]byte(fmt.Sprintf("Subject: StockReport %s\n", headers)))

	//populate the report to buffer
	err := tpl.Execute(buffer, report)
	if err != nil {
		return err
	}

	//setup fields to send email
	to := mail.NewEmail(report.User.FirstName, report.User.Email)
	subject := "Stock Report"
	from := mail.NewEmail("Stock-reporter @", os.Getenv("FROM_EMAIL"))
	htmlContent := buffer.String()
	email := mail.NewSingleEmail(from, subject, to, "None NEeded", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	resp, err := client.Send(email)
	if err != nil {
		return err
	}
	log.Printf("StatusCode: %v\n, ResponseBody: %v\n, ResponseHeader: %v\n", resp.StatusCode, resp.Body, resp.Headers)
	return nil
}
