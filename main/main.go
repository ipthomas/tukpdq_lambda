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
	pdq := initPDQ(req)
	if pdq.Used_PID == "" {
		return handle_Response("", http.StatusBadRequest, errors.New("invalid request"))
	}
	if req.QueryStringParameters["cache"] != "false" {
		if cachepat, ok := cachedpatients[pdq.Used_PID]; ok {
			log.Printf("Cached Patient found for Patient ID %s", pdq.Used_PID)
			if req.QueryStringParameters["rsptype"] == "bool" {
				return handle_Response("true", http.StatusOK, nil)
			}
			if req.QueryStringParameters["rsptype"] == "code" {
				return handle_Response("", http.StatusOK, nil)
			}
			return handle_Response(string(cachepat), http.StatusOK, nil)
		}
	}
	log.Printf("Using %s Server URL %s Patient ID %s for PDQ request", pdq.Server, pdq.Used_PID, pdq.Server_URL)
	if err := tukpdq.PDQ(&pdq); err != nil {
		return handle_Response("", pdq.StatusCode, err)
	}
	if pdq.Count < 1 {
		if req.QueryStringParameters["rsptype"] == "bool" {
			return handle_Response("false", http.StatusOK, nil)
		}
		if req.QueryStringParameters["rsptype"] == "code" {
			return handle_Response("", http.StatusNoContent, nil)
		}
		return handle_Response("No Patient Found", http.StatusOK, nil)
	}
	if req.QueryStringParameters["cache"] == "true" {
		cachedpatients[pdq.Used_PID] = pdq.Response
	}
	log.Printf("Patient ID %s is Registered", pdq.Used_PID)
	if req.QueryStringParameters["rsptype"] == "bool" {
		return handle_Response("true", http.StatusOK, nil)
	}
	if req.QueryStringParameters["rsptype"] == "code" {
		return handle_Response("", http.StatusOK, nil)
	}
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
func initPDQ(req events.APIGatewayProxyRequest) tukpdq.PDQQuery {
	var server = "pixv3"
	var nhsoid = "2.16.840.1.113883.2.1.4.1"
	var serverurl = ""
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
	switch server {
	case "pixm":
		serverurl = os.Getenv("SERVER_PIXM")
	case "pdqv3":
		serverurl = os.Getenv("SERVER_PDQV3")
	case "pixv3":
		serverurl = os.Getenv("SERVER_PIXV3")
	}
	pdq := tukpdq.PDQQuery{
		Server:     server,
		MRN_ID:     req.QueryStringParameters["mrnid"],
		MRN_OID:    req.QueryStringParameters["mrnoid"],
		NHS_ID:     req.QueryStringParameters["nhsid"],
		NHS_OID:    nhsoid,
		REG_ID:     req.QueryStringParameters["regid"],
		REG_OID:    os.Getenv("REG_OID"),
		Server_URL: serverurl,
	}
	if pdq.MRN_ID != "" && pdq.MRN_OID != "" {
		pdq.Used_PID = pdq.MRN_ID
		pdq.Used_PID_OID = pdq.MRN_OID
	} else {
		if pdq.NHS_ID != "" {
			pdq.Used_PID = pdq.NHS_ID
			pdq.Used_PID_OID = pdq.NHS_ID
		} else {
			if pdq.REG_ID != "" && pdq.REG_OID != "" {
				pdq.Used_PID = pdq.REG_ID
				pdq.Used_PID_OID = pdq.REG_OID
			}
		}
	}
	return pdq
}
