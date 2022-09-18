# tukpdq_lambda

This is an IHE PDQ Client implementation for deployment in AWS as a Lambda function. 

It uses the github.com/ipthomas/tukpdq package to perform an IHE PDQ against either :-
    An IHE PIXm compliant Server using Fhir
    An IHE PIXv3 compliant Server using SOAP
    An IHE PDQv3 compliant Server using SOAP

Required AWS Environment Variables are:
    Key                                         (Example Values)
    NHS_OID	                                    2.16.840.1.113883.2.1.4.1
    REG_OID	                                    2.16.840.1.113883.2.1.3.31.2.1.1
    SERVER_DEFAULT	                            SERVER_PIXV3
    SERVER_PDQV3	                            http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PDQSupplier
    SERVER_PIXM	                                http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIXFhir/r4/Patient
    SERVER_PIXV3	                            http://spirit-test-01.tianispirit.co.uk:8081/SpiritPIX/PIXManager

It takes an AWS APIProxyRequest input and parses the requestParams:
    If queryParam server is set, the value will overide the SERVER_DEFAULT Env var: Valid values are:
        server=pixm
        server=pixv3
        server=pdqv3
   
    The PDQ client supports using either a MRN ID along with the MRN OID or a NHS ID or a Reg ID for the Patient ID to use in the query. If nhsoid or regoid queryParam values are provided they will override the env vars NHS_OID, REG_OID. 
    The queryParams available for specifying the patient id are:
        mrnid= and mrnoid=
        nhsid= (and optional nhsoid=)
        regid= (and optional regoid=)
    
    Patients are cached by default. To bypass the cache use the queryParam cache=false. The cache is valid as long as the state is maintained by the AWS lambda function.

The queryParam rsptype specifies the required response to the PDQ. 
    If queryParam rsptype="bool", the returned response is either 'true' or 'false'.
    If queryParam rsptype="code", the http header 'statuscode' is set to either 200 if patient is found or 204 if not. The response body is empty.
    The default returnes either the patient details or 'No Patient Found"