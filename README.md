# tukpdq_lambda

This is an implementation of IHE PDQ Clients for deployment in AWS as a Lambda function. 

The PDQ is performed against either :-
    An IHE PIXm compliant Server using Fhir
    An IHE PIXv3 compliant Server using SOAP
    An IHE PDQv3 compliant Server using SOAP

AWS Environment Variables are:
    Key                                         (Example Values)
    NHS_OID                                     The NHS Default will be used if non provided
    REG_OID	                                    2.16.840.1.113883.2.1.3.31.2.1.1
    PDQ_SERVER_TYPE	                            pdqv3
    PDQ_SERVER_URL	                            http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PDQSupplier
    PATIENT_CACHE                               true

It takes an AWS APIProxyRequest input with the nhsid specified as a requestParam 
    Eg. https://k7mmeyp191.execute-api.eu-west-1.amazonaws.com/beta/ping?nhsid=2222222222