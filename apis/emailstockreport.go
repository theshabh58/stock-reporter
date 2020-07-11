package apis

import (
	"bytes"
	"context"
	"encoding/json"
	"html/template"
	"log"
	"net/smtp"
	"os"
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

	//setup auth
	emailAPIUsr := os.Getenv("EMAIL_API_USERNAME")
	emailAPIPswd := os.Getenv("EMAIL_API_PASSWORD")
	emailAPIHst := os.Getenv("EMAIL_API_HOST")
	emailAuth = smtp.PlainAuth("Stock-Reporter-Api", emailAPIUsr, emailAPIPswd, emailAPIHst)
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

	//populate the report to buffer
	err := tpl.Execute(buffer, report.StockReport)
	if err != nil {
		return err
	}

	//setup email
	emailParams := EmailParameters{
		To:      []string{report.User.Email},
		From:    "stockreporter-28c938@inbox.mailtrap.io",
		Subject: "Stock-Report",
		Message: buffer.Bytes(),
	}

	err = smtp.SendMail("smtp.mailtrap.io:2525", emailAuth, emailParams.From, emailParams.To, emailParams.Message)

	if err != nil {
		return err
	}
	return nil
}
