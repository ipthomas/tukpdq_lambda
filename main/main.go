package main

import (
	"os"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ipthomas/tukcnst"
	"github.com/ipthomas/tukpdq"
)

func main() {
	lambda.Start(Handle_Request)
}

// Set AWS Env PDQ_SERVER_URL to the WSE for the IHE Complaint PDQ Server you wish to use
//
//	EG
//		PDQv3 server wse - http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PDQSupplier
//
// or	PIXv3 server wse - http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PIXManager
//
// or 	PIXm server wse  - http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient
//
// Set AWS Env PDQ_SERVER_TYPE to specify the PDQ server type.
//
//	Valid types are
//		PDQv3 SOAP server - pdqv3
//
// or 	PIXv3 SOAP server - pixv3
//
// or 	PIXm  FHIR server - pixm
//
// Set AWS Env Reg_OID to the regional oid
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var err error
	patcache, _ := strconv.ParseBool(os.Getenv(tukcnst.AWS_ENV_PATIENT_CACHE))
	pdq := tukpdq.PDQQuery{
		Server:     os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_TYPE),
		MRN_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_ID],
		MRN_OID:    req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_OID],
		NHS_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_ID],
		NHS_OID:    os.Getenv(tukcnst.AWS_ENV_REG_OID),
		REG_ID:     req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_ID],
		REG_OID:    os.Getenv(tukcnst.AWS_ENV_REG_OID),
		Server_URL: os.Getenv(tukcnst.AWS_ENV_PDQ_SERVER_URL),
		RspType:    os.Getenv(tukcnst.AWS_ENV_RESPONSE_TYPE),
		Cache:      patcache,
		Timeout:    5,
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_OID] != "" {
		pdq.NHS_OID = req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_OID]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_OID] != "" {
		pdq.REG_OID = req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_OID]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE] != "" && req.QueryStringParameters["pdqserverurl"] != "" {
		pdq.Server = req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE]
		pdq.Server_URL = req.QueryStringParameters["pdqserverurl"]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_RESPONSE_TYPE] != "" {
		pdq.RspType = req.QueryStringParameters[tukcnst.QUERY_PARAM_RESPONSE_TYPE]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_CACHE] != "" {
		pdqcache, _ := strconv.ParseBool(os.Getenv(tukcnst.AWS_ENV_PATIENT_CACHE))
		pdq.Cache = pdqcache
	}

	err = tukpdq.New_Transaction(&pdq)
	resp := events.APIGatewayProxyResponse{
		StatusCode: pdq.StatusCode,
		Body:       string(pdq.Response),
	}
	return &resp, err
}
