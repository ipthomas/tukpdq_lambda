package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ipthomas/tukpdq"
)

var cachedpatients = make(map[string][]byte)

func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var st = time.Now()
	var et = st
	var server = strings.ToLower(strings.TrimPrefix(os.Getenv("SERVER_DEFAULT"), "SERVER_"))
	var nhsoid = os.Getenv("NHS_OID")
	if req.QueryStringParameters["server"] != "" {
		server = strings.ToLower(req.QueryStringParameters["server"])
	}
	if req.QueryStringParameters["nhsoid"] != "" {
		nhsoid = "2.16.840.1.113883.2.1.4.1"
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
	usedPID := getUsedPID(pdq)
	if usedPID == "" {
		return handle_Response("", http.StatusBadRequest, errors.New("invalid request"))
	}
	if req.QueryStringParameters["cache"] == "true" {
		if isregistered, ok := cachedpatients[usedPID]; ok {
			log.Printf("Cached ID %s found for registered patient", usedPID)
			return handle_Response(string(isregistered), http.StatusOK, nil)
		}
	}
	log.Printf("Using %s server for PDQ request", pdq.Server)
	if err := tukpdq.PDQ(&pdq); err != nil {
		return handle_Response("", pdq.StatusCode, err)
	}
	if pdq.Count < 1 {
		return handle_Response("No Patient Found", http.StatusNoContent, nil)
	}
	cachedpatients[pdq.Used_PID] = pdq.Response
	log.Printf("Patient ID %s is Registered", pdq.Used_PID)
	et = time.Now()
	log.Printf("Response time %v ms", et.Sub(st).Milliseconds())
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
func getUsedPID(i tukpdq.PDQQuery) string {
	if i.MRN_ID != "" && i.MRN_OID != "" {
		return i.MRN_ID
	} else {
		if i.NHS_ID != "" {
			return i.NHS_ID
		} else {
			if i.REG_ID != "" && i.REG_OID != "" {
				return i.REG_ID
			}
		}
	}
	return ""
}
