package auth

import (
	"encoding/xml"
	"fmt"
)

type ErrEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    ErrBody  `xml:"Body"`
}

type ErrBody struct {
	Fault Fault `xml:"Fault"`
}

type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
	Detail      Detail `xml:"detail"`
}

type Detail struct {
	ExceptionName string `xml:"exceptionName"`
	Hostname      string `xml:"hostname"`
}

type ErrDetails struct {
	FaultCode string
}

func ParseResponseError(resBytes []byte) (Fault, error) {
	var envelope ErrEnvelope

	err := xml.Unmarshal(resBytes, &envelope)
	if err != nil {
		return Fault{}, fmt.Errorf("error unmarshaling error response: %v", err)
	}

	return envelope.Body.Fault, err
}
