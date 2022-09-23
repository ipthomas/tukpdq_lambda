package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/ipthomas/tukint"
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
// Set AWS Env PDQ_SERVER_DEFAULT to specify the PDQ server type.
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
	return tukint.NewPDQ(req)
}
