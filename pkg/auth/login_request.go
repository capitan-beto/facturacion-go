package auth

import (
	"encoding/xml"
	"fmt"

	"cmd/api/main.go/pkg/utils"
)

type Envelope struct {
	XMLName   xml.Name  `xml:"soapenv:Envelope"`
	XMLNSs    string    `xml:"xmlns:soapenv,attr"`
	XMLNSwsaa string    `xml:"xmlns:wsaa,attr"`
	Header    string    `xml:"soapenv:Header"`
	Body      *WSAABody `xml:"soapenv:Body"`
}

type WSAABody struct {
	LoginCms *LoginCms `xml:"wsaa:loginCms"`
}

type LoginCms struct {
	Cms xml.CharData `xml:"wsaa:in0"`
}

// CreateLoginRequest creates an xml request appending the loginCMS
func CreateLoginRequest(loginCMS []byte) ([]byte, error) {
	soapenv := "http://schemas.xmlsoap.org/soap/envelope/"
	wsaa := "http://wsaa.view.sua.dvadac.desein.afip.gov"

	encodedCms, err := FormatCMS(loginCMS)
	if err != nil {
		return nil, fmt.Errorf("error reading cms: %v", err)
	}

	newCmsField := LoginCms{Cms: []byte(encodedCms)}
	newWSAABody := WSAABody{LoginCms: &newCmsField}
	newReq := Envelope{XMLNSs: soapenv, XMLNSwsaa: wsaa, Body: &newWSAABody}

	out, _ := xml.MarshalIndent(newReq, " ", "  ")
	res := utils.EscapeXML(string(out))

	return []byte(res), err
}

func FormatCMS(cms []byte) (string, error) {
	res := string(cms)
	return res[20 : len(res)-19], nil
}
