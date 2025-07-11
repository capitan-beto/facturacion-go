package auth

import (
	"testing"
)

var testData = `<?xml version="1.0" encoding="UTF-8"?>
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
        <soapenv:Body>
            <soapenv:Fault>
                <faultcode xmlns:ns1="http://xml.apache.org/axis/">ns1:coe.alreadyAuthenticated</faultcode>
                <faultstring>El CEE ya posee un TA valido para el acceso al WSN solicitado</faultstring>
                <detail>
                    <ns2:exceptionName xmlns:ns2="http://xml.apache.org/axis/">gov.afip.desein.dvadac.sua.view.wsaa.LoginFault</ns2:exceptionName>
                    <ns3:hostname xmlns:ns3="http://xml.apache.org/axis/">wsaaext1.homo.afip.gov.ar</ns3:hostname>
                </detail>
        </soapenv:Fault>
    </soapenv:Body>
</soapenv:Envelope>`

func TestParseResponseError(t *testing.T) {

	res, err := ParseResponseError([]byte(testData))
	if err != nil {
		t.Fatalf("error! expected nil err, got %v", err)
	}

	expected := "El CEE ya posee un TA valido para el acceso al WSN solicitado"

	if res.FaultString != expected {
		t.Fatalf("error! expected %v, got %v", expected, res.FaultString)
	}

}
