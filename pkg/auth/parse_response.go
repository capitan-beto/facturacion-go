package auth

import (
	"encoding/xml"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type MyResEnvelope struct {
	XMLName xml.Name
	Body    Body
}

type Body struct {
	XMLName          xml.Name
	LoginCmsResponse LoginCmsResponse `xml:"loginCmsResponse"`
}

type LoginCmsResponse struct {
	LoginCmsReturn string `xml:"loginCmsReturn"`
}

type LoginTicketResponse struct {
	XMLName     xml.Name    `xml:"loginTicketResponse"`
	Version     string      `xml:"version,attr"`
	Header      HeaderRes   `xml:"header"`
	Credentials Credentials `xml:"credentials"`
}

type HeaderRes struct {
	Source         string `xml:"source"`
	Destination    string `xml:"destination"`
	UniqueId       string `xml:"uniqueId"`
	GenerationTime string `xml:"generationTime"`
	ExpirationTime string `xml:"expirationTime"`
}

type Credentials struct {
	Token string `xml:"token"`
	Sign  string `xml:"sign"`
}

func ParseResponse(resBytes []byte) (LoginTicketResponse, error) {
	resEnvelope := &MyResEnvelope{}

	err := xml.Unmarshal(resBytes, &resEnvelope)
	if err != nil {
		log.Error(err)
		return LoginTicketResponse{}, fmt.Errorf("error unmarshaling response: %v", err)
	}

	var ltr LoginTicketResponse
	err = xml.Unmarshal([]byte(resEnvelope.Body.LoginCmsResponse.LoginCmsReturn), &ltr)
	if err != nil {
		return LoginTicketResponse{}, err
	}

	return ltr, err
}
