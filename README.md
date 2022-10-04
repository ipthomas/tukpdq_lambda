# tukpdq_lambda

This is an implementation of IHE PDQ Clients (PIXv3, PDQv3 and PIXm) for deployment in AWS as a Lambda function. It also supports querying the CGL service (drug and substance use) and reutrns the CGL 'user' patient if registered with CGL. Note the CGL PDQ response 'user' contains both demographics and CGL content.

The PDQ is performed against either :-
    An IHE PIXm compliant Server using Fhir/json
    An IHE PIXv3 compliant Server using SOAP/xml
    An IHE PDQv3 compliant Server using SOAP/xml
    CGL Server using REST/json

AWS Environment Variables are:
                                                             (Example Values)
    NHS_OID                                     2.16.840.1.113883.2.1.4.1 (The NHS Default will be used if non provided)
    REG_OID	                                    2.16.840.1.113883.2.1.3.31.2.1.1 (Must be set or provided in query)
    PDQ_SERVER_TYPE	                            pdqv3 (Must be set as env var or provided in query)
    PDQ_SERVER_URL	                            http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PDQSupplier (Must be set as env var or provided in query)
    PATIENT_CACHE                               true (Default is false). Query param cache= will overide Env var
    CGL_API_KEY                                 FNhb#OhxWiEiMdf+@6085k5Zmt (Optional unless PDQ_SERVER_TYPE=cgl or you want to perform an additional query against the CGL server along with the IHE PDQ query
    CGL_SERVER_URL                              https://public-api.criisdev.org.uk/api/v1/user?NHS_number= (Optional unless PDQ_SERVER_TYPE = cgl or the additional PDQ against the CGL server is required)

Example AWS API G/W request:
https://k6mmeyp391.execute-api.eu-west-1.amazonaws.com/beta/ping?nhsid=6072406157&cache=false&pdqserver=pdqv3&_include=cgl

The build folder contains an AWS Lambda build (main.zip)
To build for AWS Lambda deployment
    GOOS=linux go build -o build/main main/main.go
    zip -jrm build/main.zip build/main 