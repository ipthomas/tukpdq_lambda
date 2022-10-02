# tukpdq_lambda

This is an implementation of IHE PDQ Clients (PIXv3, PDQv3 and PIXm) for deployment in AWS as a Lambda function. It also supports querying the CGL service (drug and substance use) and reutrns the CGL 'user' patient if registered with CGL. The CGL 'user' contains both demographics and CGL content.

The PDQ is performed against either :-
    An IHE PIXm compliant Server using Fhir
    An IHE PIXv3 compliant Server using SOAP
    An IHE PDQv3 compliant Server using SOAP

AWS Environment Variables are:
    Key                                         (Example Values)
    NHS_OID                                     2.16.840.1.113883.2.1.4.1 (The NHS Default will be used if non provided)
    REG_OID	                                    2.16.840.1.113883.2.1.3.31.2.1.1 (Must be set or provided in query)
    PDQ_SERVER_TYPE	                            pdqv3 (Must be set or provided in query)
    PDQ_SERVER_URL	                            http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PDQSupplier (Must be set or provided in query)
    PATIENT_CACHE                               true (Default is false)
    CGL_API_KEY                                 FNhb#OhxWiEiMdf+@6085k5Zmt (Optional unless PDQ_SERVER_TYPE = cgl. If present, an additonal query is made against the CGL_SERVER_URL and the PDQ.CGL_User is populated if the patient is also registered with CGL
    CGL_SERVER_URL                              https://public-api.criisdev.org.uk/api/v1/user?NHS_number= (Optional unless PDQ_SERVER_TYPE = cgl or the additional PDQ against the CGL server is required)

It takes an AWS APIProxyRequest input with a patient nhs id specified in the requestParams 
    Eg. https://k7mmeyp191.execute-api.eu-west-1.amazonaws.com/beta/ping?nhsid=2222222222