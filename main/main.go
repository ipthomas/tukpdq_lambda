package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ipthomas/tukpdq"
)

var cachedpatients = make(map[string][]byte)

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var server = "pixv3"
	var nhsoid = "2.16.840.1.113883.2.1.4.1"
	if os.Getenv("NHS_OID") != "" {
		nhsoid = os.Getenv("NHS_OID")
	}
	if req.QueryStringParameters["nhsoid"] != "" {
		nhsoid = req.QueryStringParameters["nhsoid"]
	}
	if os.Getenv("SERVER_DEFAULT") != "" {
		server = strings.ToLower(strings.TrimPrefix(os.Getenv("SERVER_DEFAULT"), "SERVER_"))
	}
	if req.QueryStringParameters["server"] != "" {
		server = strings.ToLower(req.QueryStringParameters["server"])
	}
	pdq := tukpdq.PDQQuery{
		Server:     server,
		MRN_ID:     req.QueryStringParameters["mrnid"],
		MRN_OID:    req.QueryStringParameters["mrnoid"],
		NHS_ID:     req.QueryStringParameters["nhsid"],
		NHS_OID:    nhsoid,
		REG_ID:     req.QueryStringParameters["regid"],
		REG_OID:    os.Getenv("REG_OID"),
		Server_URL: getServerURL(server),
	}
	setQueryPID(&pdq)
	if pdq.Used_PID == "" {
		return handle_Response("", http.StatusBadRequest, errors.New("invalid request"))
	}
	if req.QueryStringParameters["cache"] == "true" {
		if cachepat, ok := cachedpatients[pdq.Used_PID]; ok {
			log.Printf("Cached ID %s found for registered patient", pdq.Used_PID)
			return handle_Response(string(cachepat), http.StatusOK, nil)
		}
	}
	log.Printf("Using %s server for PDQ request", pdq.Server)
	if err := tukpdq.PDQ(&pdq); err != nil {
		return handle_Response("", pdq.StatusCode, err)
	}
	if pdq.Count < 1 {
		return handle_Response("No Patient Found", http.StatusNoContent, nil)
	}
	if req.QueryStringParameters["cache"] == "true" {
		cachedpatients[pdq.Used_PID] = pdq.Response
	}
	log.Printf("Patient ID %s is Registered", pdq.Used_PID)
	return handle_Response(string(pdq.Response), http.StatusOK, nil)
}
func handle_Response(body string, status int, err error) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{MultiValueHeaders: map[string][]string{}, IsBase64Encoded: false}
	if err == nil {
		resp.StatusCode = status
		resp.Body = body
	} else {
		resp.StatusCode = status
		resp.Body = err.Error()
		log.Println(err.Error())
	}
	return &resp, err
}
func getServerURL(server string) string {
	switch server {
	case "pixm":
		return os.Getenv("SERVER_PIXM")
	case "pdqv3":
		return os.Getenv("SERVER_PDQV3")
	}
	return os.Getenv("SERVER_PIXV3")
}
func setQueryPID(i *tukpdq.PDQQuery) {
	if i.MRN_ID != "" && i.MRN_OID != "" {
		i.Used_PID = i.MRN_ID
		i.Used_PID_OID = i.MRN_OID
	} else {
		if i.NHS_ID != "" {
			i.Used_PID = i.NHS_ID
			i.Used_PID_OID = i.NHS_ID
		} else {
			if i.REG_ID != "" && i.REG_OID != "" {
				i.Used_PID = i.REG_ID
				i.Used_PID_OID = i.REG_OID
			}
		}
	}
}
