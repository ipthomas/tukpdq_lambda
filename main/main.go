package main

import (
	"encoding/json"
	"log"
	"net/http"
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
// or 	CGL   HTTP server - cgl
//
// Set AWS Env Reg_OID to the regional oid
//
// A PDQ against any of the 3 IHE PDQ server types can also include the results of a query against the CGL service if the CGL_API_KEY and CGL_SERVER_URL are set
// To perform just a query against the CGL service, set PDQ_SERVER_TYPE=cgl
func Handle_Request(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	patcache, _ := strconv.ParseBool(os.Getenv(tukcnst.ENV_PATIENT_CACHE))
	pdq := tukpdq.PDQQuery{
		Server_Mode:   os.Getenv(tukcnst.ENV_PDQ_SERVER_TYPE),
		CGL_X_Api_Key: os.Getenv(tukcnst.ENV_CGL_X_API_KEY),
		MRN_ID:        req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_ID],
		MRN_OID:       req.QueryStringParameters[tukcnst.QUERY_PARAM_MRN_OID],
		NHS_ID:        req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_ID],
		NHS_OID:       os.Getenv(tukcnst.ENV_NHS_OID),
		REG_ID:        req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_ID],
		REG_OID:       os.Getenv(tukcnst.ENV_REG_OID),
		Server_URL:    os.Getenv(tukcnst.ENV_PDQ_SERVER_URL),
		Cache:         patcache,
		Timeout:       5,
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_OID] != "" {
		pdq.NHS_OID = req.QueryStringParameters[tukcnst.QUERY_PARAM_NHS_OID]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_OID] != "" {
		pdq.REG_OID = req.QueryStringParameters[tukcnst.QUERY_PARAM_REG_OID]
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE] != "" {
		log.Printf("Setting Server type to %s", req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE])
		srvurl := getPDQServerURL(req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE])
		if srvurl != "" {
			pdq.Server_URL = srvurl
			log.Printf("Set Server URL to %s", pdq.Server_URL)
			pdq.Server_Mode = req.QueryStringParameters[tukcnst.QUERY_PARAM_PDQ_SERVER_TYPE]
			log.Printf("Set Server type to %s", pdq.Server_Mode)
		}
	}
	if req.QueryStringParameters[tukcnst.QUERY_PARAM_CACHE] != "" {
		pdqcache, _ := strconv.ParseBool(req.QueryStringParameters[tukcnst.QUERY_PARAM_CACHE])
		pdq.Cache = pdqcache
	}

	if err := tukpdq.New_Transaction(&pdq); err != nil {
		log.Println(err.Error())
	}

	if pdq.Server_Mode != tukcnst.PDQ_SERVER_TYPE_CGL && pdq.CGL_X_Api_Key != "" && req.QueryStringParameters[tukcnst.QUERY_PARAM_INCLUDE] == tukcnst.PDQ_SERVER_TYPE_CGL {
		log.Println("Performing additional query against CGL service")
		cglpdq := tukpdq.PDQQuery{
			Server_Mode:   tukcnst.PDQ_SERVER_TYPE_CGL,
			CGL_X_Api_Key: pdq.CGL_X_Api_Key,
			NHS_ID:        pdq.NHS_ID,
			NHS_OID:       tukcnst.NHS_OID_DEFAULT,
			REG_OID:       pdq.REG_OID,
			Server_URL:    getPDQServerURL(tukcnst.PDQ_SERVER_TYPE_CGL),
		}
		err := tukpdq.New_Transaction(&cglpdq)
		if err != nil {
			log.Println(err.Error())
		}
		pdq.CGLUserResponse = cglpdq.CGLUserResponse
	}
	var b []byte
	b, _ = json.MarshalIndent(pdq, "", "  ")
	apiResp := events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(b),
	}
	return &apiResp, nil
}
func getPDQServerURL(srv string) string {
	log.Printf("Selecting %s Server URL", srv)
	srvurl := ""
	switch srv {
	case tukcnst.PDQ_SERVER_TYPE_CGL:
		srvurl = os.Getenv(tukcnst.ENV_CGL_SERVER_URL)
	case tukcnst.PDQ_SERVER_TYPE_IHE_PDQV3:
		srvurl = os.Getenv(tukcnst.ENV_IHE_PDQV3_SERVER_URL)
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXV3:
		srvurl = os.Getenv(tukcnst.ENV_IHE_PIXV3_SERVER_URL)
	case tukcnst.PDQ_SERVER_TYPE_IHE_PIXM:
		srvurl = os.Getenv(tukcnst.ENV_IHE_PIXM_SERVER_URL)
	}
	log.Printf("Selected %s server URL %s", srv, srvurl)
	return srvurl
}
