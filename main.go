package main

/*
* Version 0.0.1
* Compatible on ANY browser
 */

/*** WORKFLOW ***/
/*
* 1- Get API request from API gateway
* 2- Determine if request if POST, ignore every other request
* 3- Parse body from multiform url
* 4- Parse request body and validate body, reject if violated
* 5- If a request contains, contact_email, use contact_email as the 'reply to' option in the email body
*    otherwise, use owner's email as the 'reply to'
* 6- Pass message content (including optional phone number and email) to beautify
* 7- Set Header status to '200 OK' if everything is successful or 403 if something went wrong
* 8- Return json with 'Message sent' body
 */

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"regexp"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	//checkmail "github.com/warrensbox/falcon-form/lib"
	lib "github.com/warrensbox/falcon-form/lib"
)

var (
	HTTPMethodNotSupported = errors.New("no name was provided in the HTTP body")
)

type Response struct {
	StatusCode int `json:"statusCode"`
	//Headers    map[string]string `json:"headers"`
	Body string `json:"body"`
}

type Message struct {
	OwnerEmail     string `json:"owner_email"`
	ContactEmail   string `json:"contact_email,omitempty"`
	ContactName    string `json:"contact_name,omitempty"`
	ContactPhone   string `json:"contact_phone,omitempty"`
	MessageContent string `json:"message_content"`
}

//test comment
const (
	URLPARAM     = "/form"
	CONTACTEMAIL = "contact_email"
	CONTACTPHONE = "contact_phone"
	CONTACTNAME  = "contact_name"
	MSGCONTENT   = "message_content"
	DEFAULT      = "http://example.com"
)

var (
	ownerEmail   string
	contactEmail string
	contactPhone string
	msgContent   string
	contactName  string
	statusCode   int
)

//HandleRequest incoming request
func HandleRequest(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	redirect := make(map[string]string)
	statusCode = 200

	redirect["Location"] = DEFAULT
	redirect["Access-Control-Allow-Origin"] = "*"
	redirect["Access-Control-Allow-Headers"] = "*"

	if request.HTTPMethod == "GET" {
		fmt.Printf("GET METHOD\n")
		statusCode = 301
		return events.APIGatewayProxyResponse{Headers: redirect, StatusCode: statusCode}, nil
	} else if request.HTTPMethod == "POST" {

		fmt.Printf("POST METHOD\n")

		fmt.Println(request)

		parameter := request.PathParameters["proxy"]
		fmt.Println(request.PathParameters)
		fmt.Println(request.Path)

		re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

		if request.Path == URLPARAM {
			fmt.Println("It's json")
			fmt.Println(parameter)

			var content Message

			data := []byte(request.Body)
			json.Unmarshal(data, &content)

			ownerEmail = content.OwnerEmail
			contactEmail = content.ContactEmail
			contactName = content.ContactName
			contactPhone = content.ContactPhone
			msgContent = content.MessageContent

			statusCode = 200

		} else if re.MatchString(parameter) {
			fmt.Println("It's a form")
			errFormat := lib.ValidateFormat(parameter)
			lib.ErrorExit("Unable to validate email format", errFormat)

			errDomain := lib.ValidateHost(parameter)
			lib.ErrorExit("Unable to validate email host", errDomain)

			ownerEmail = parameter

			urlSample, errURL := url.Parse("http://test.com?" + request.Body)
			lib.ErrorExit("Unable to parse url", errURL)

			fmt.Printf("url sample %q \n", urlSample)
			content, _ := url.ParseQuery(urlSample.RawQuery)

			if len(content[CONTACTEMAIL]) > 0 {
				fmt.Printf("Contact email %q \n", content[CONTACTEMAIL][0])
				contactEmail = content[CONTACTEMAIL][0]
			} else {
				fmt.Printf("no email")
				contactEmail = ""
			}

			if len(content[CONTACTPHONE]) > 0 {
				fmt.Printf("phone %q \n", content[CONTACTPHONE][0])
				contactPhone = content[CONTACTPHONE][0]
			} else {
				fmt.Printf("no phone")
				contactPhone = ""
			}

			if len(content[CONTACTNAME]) > 0 {
				fmt.Printf("name %q \n", content[CONTACTNAME][0])
				contactName = content[CONTACTNAME][0]
			} else {
				fmt.Printf("no name")
				contactName = ""
			}

			if len(content[MSGCONTENT]) > 0 {
				fmt.Printf("msg %q \n", content[MSGCONTENT][0])
				msgContent = content[MSGCONTENT][0]
			} else {
				fmt.Printf("no msg")
				msgContent = ""
			}

			statusCode = 301

		} else {
			return events.APIGatewayProxyResponse{}, HTTPMethodNotSupported
		}

		msgInfo := lib.SendEmail(ownerEmail, contactEmail, contactPhone, contactName, msgContent)
		fmt.Println(msgInfo)

		if request.Headers["origin"] != "" {
			redirect["Location"] = request.Headers["origin"]
		}

		return events.APIGatewayProxyResponse{Headers: redirect, StatusCode: statusCode}, nil
	} else {
		fmt.Printf("NEITHER\n")
		return events.APIGatewayProxyResponse{}, HTTPMethodNotSupported
	}

}

func main() {
	lambda.Start(HandleRequest)
}
