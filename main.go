package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ssm "github.com/aws/aws-sdk-go/service/ssm"
	"github.com/matcornic/hermes"
	gophermail "gopkg.in/jpoehls/gophermail.v0"
)

// ErrNameNotProvided is thrown when a name is not provided

var (
	HTTPMethodNotSupported = errors.New("no name was provided in the HTTP body")
)

type Response struct {
	StatusCode int               `json:"statusCode"`
	Headers    map[string]string `json:"headers"`
	Body       string            `json:"body"`
}

type Thing struct {
	Field1 int    `json:"field1"`
	Field2 string `json:"field2"`
}

type Message struct {
	OwnerEmail     string `json:"owner_email"`
	ContactEmail   string `json:"contact_email,omitempty"`
	MessageContent string `json:"message_content"`
}

func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	fmt.Printf("Body size = %d. \n", len(request.Body))
	fmt.Println("Headers:")
	var msgContent Message
	data := []byte(request.Body)
	json.Unmarshal(data, &msgContent)

	fmt.Println("Printing data")
	fmt.Println(msgContent)
	fmt.Println("Printing data")

	t := Thing{45, "world"}
	j, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}

	session := session.Must(session.NewSession())
	svc := ssm.New(session)

	pass := &ssm.GetParameterInput{
		Name:           aws.String("SMTP_PASS"),
		WithDecryption: aws.Bool(true),
	}

	respPass, err := svc.GetParameter(pass)
	errorExit("GetParameters", err)

	SMTPPASS := *respPass.Parameter.Value
	fmt.Println(SMTPPASS)

	user := &ssm.GetParameterInput{
		Name:           aws.String("SMTP_USER"),
		WithDecryption: aws.Bool(true),
	}

	respUser, err := svc.GetParameter(user)

	SMTPUSER := *respUser.Parameter.Value
	fmt.Println(SMTPUSER)

	smtpEmail := &ssm.GetParameterInput{
		Name:           aws.String("SMTP_EMAIL"),
		WithDecryption: aws.Bool(true),
	}

	respEmail, err := svc.GetParameter(smtpEmail)

	SMTPEMAIL := *respEmail.Parameter.Value
	fmt.Println(SMTPEMAIL)

	smtpPort := &ssm.GetParameterInput{
		Name:           aws.String("SMTP_PORT"),
		WithDecryption: aws.Bool(true),
	}

	respPort, err := svc.GetParameter(smtpPort)

	SMTPPORT := *respPort.Parameter.Value
	fmt.Println(SMTPPORT)

	url := "https://resumex.io"
	name := "Becks"

	h := hermes.Hermes{
		// Optional Theme
		Theme: new(hermes.Flat),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name: "Warren from resumex.io",
			Link: url,
			// Optional product logo
			//Logo:      imageHeader,
			Copyright: "Ⓒ 2017 Keplerbox LLC - Crafted with ❤ in San Francisco, California",
		},
	}

	emailcontent := WelcomeEmail(name, url)

	emailBody, errBody := h.GenerateHTML(emailcontent)
	if errBody != nil {
		fmt.Println(errBody)
	}

	// emailText, errText := h.GeneratePlainText(emailcontent)
	// if errText != nil {
	// 	fmt.Println(errText)
	// }

	// e := em.NewEmail()
	// e.From = "warren.veerasingam@gmail.com"
	// e.To = []string{"umesh.veerasingam@gmail.com"}
	// e.Subject = "Ready to Kickass?"
	// e.Text = []byte(emailText)
	// e.HTML = []byte(emailBody)

	// auth := smtp.PlainAuth("", SMTPUSER, SMTPPASS, SMTPEMAIL)
	// errEmail := e.Send(SMTPPORT, auth)
	// if errEmail != nil {
	// 	fmt.Println(errEmail)
	// }

	fromEmail := "support@resumex.io"
	to_email := "umesh.veerasingam@gmail.com"

	var msg gophermail.Message
	msg.SetFrom(fromEmail)
	msg.AddTo(to_email)
	msg.SetReplyTo(to_email)
	msg.Subject = "Would you take a look at my resume?"
	msg.HTMLBody = emailBody

	auth2 := smtp.PlainAuth("", SMTPUSER, SMTPPASS, SMTPEMAIL)
	errEmail2 := gophermail.SendMail(SMTPPORT, auth2, &msg)
	if errEmail2 != nil {
		fmt.Println(errEmail2)
	}

	for key, value := range request.Headers {
		fmt.Printf("  %s: %s\n", key, value)
	}
	if request.HTTPMethod == "GET" {
		fmt.Printf("GET METHOD\n")
		return events.APIGatewayProxyResponse{Body: string(j), StatusCode: 200}, nil
	} else if request.HTTPMethod == "POST" {
		fmt.Printf("POST METHOD\n")
		fmt.Printf(request.Body)
		return events.APIGatewayProxyResponse{Body: request.Body, StatusCode: 200}, nil
	} else {
		fmt.Printf("NEITHER\n")
		return events.APIGatewayProxyResponse{}, HTTPMethodNotSupported
	}
}

func main() {

	lambda.Start(HandleRequest)
}

func errorExit(msg string, e error) {
	if e != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s, %v\n", msg, e)
		os.Exit(1)
	}
}

func WelcomeEmail(name string, url string) hermes.Email {
	return hermes.Email{
		Body: hermes.Body{
			Name: name,
			Intros: []string{
				"Welcome to ResumeX! We're very excited to have you on board.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "To get started with ResumeX, please click here:",
					Button: hermes.Button{
						Text: "Start Working On My Resume",
						Link: url + "/user",
					},
				},
			},
			Greeting: "Hello",
			Outros: []string{
				"Need help, or have questions? Just send an email to support@resumex.io, we'd love to help.",
			},
		},
	}
}
