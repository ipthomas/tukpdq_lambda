package main

import (
	"log"
	"net/http"
	"time"
	"tukdsub"
	"tukint"

	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukpixm"
	"github.com/ipthomas/tukutil"
)

func main() {
	var st = time.Now()
	var et = time.Now()
	var pid = "1111111111"
	var pidoid = "2.16.840.1.113883.2.1.4.1"
	var regoid = "2.16.840.1.113883.2.1.3.31.2.1.1"
	var serverurl = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient"
	log.Printf("Patient Registered %v", PIXm_IsPatientRegistered(pid, pidoid, regoid, serverurl))
	et = time.Now()
	dur := et.Sub(st).Milliseconds()
	log.Printf("Duration %vms", dur)
	st = time.Now()
	pats := PatientQuery(pid, pidoid, regoid, serverurl)
	tukutil.Log(pats)
	et = time.Now()
	dur = et.Sub(st).Milliseconds()
	log.Printf("Duration %vms", dur)
}
func PIXm_IsPatientRegistered(pid string, pidoid string, regoid string, serverurl string) bool {
	pdq := tukpixm.PDQQuery{
		Server:     tukcnst.PIXm,
		NHS_ID:     pid,
		NHS_OID:    pidoid,
		REG_OID:    regoid,
		Server_URL: serverurl,
	}
	tukpixm.PDQ(&pdq)
	return pdq.Count > 0
}
func PatientQuery(pid string, pidoid string, regoid string, serverurl string) []tukpixm.PIXPatient {
	pdq := tukpixm.PDQQuery{
		Server:     tukcnst.PIXm,
		NHS_ID:     pid,
		NHS_OID:    pidoid,
		REG_OID:    regoid,
		Server_URL: serverurl,
	}
	tukpixm.PDQ(&pdq)
	return pdq.Patients
}
func main() {
	lambda.Start(Handle_Request)
}
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var st = time.Now()
	var et = time.Now()
	var pidoid = req.QUERY_PARAM_NHS_ID
	var regoid = "2.16.840.1.113883.2.1.3.31.2.1.1"
	var serverurl = "http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient"

	dsubmsg := tukint.EventMessage{Message: req.Body}
	return handle_Response(dsubmsg.NewDSUBBrokerEvent())
}
func handle_Response(err error) (*events.APIGatewayProxyResponse, error) {
	resp := events.APIGatewayProxyResponse{Headers: tukint.SOAP_XML_Content_Type_EventHeaders, MultiValueHeaders: map[string][]string{}, IsBase64Encoded: false}
	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Body = tukdsub.DSUB_ACK_TEMPLATE
	} else {
		resp.StatusCode = http.StatusInternalServerError
		resp.Body = err.Error()
		log.Println(err.Error())
	}
	return &resp, err
}
