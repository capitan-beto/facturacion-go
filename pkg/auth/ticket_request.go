package auth

import (
	"encoding/xml"
	"fmt"
	"math/rand"
	"time"

	"github.com/beevik/ntp"
)

// Body of ticket request.
type TicketRequest struct {
	XMLName xml.Name `xml:"loginTicketRequest"`
	Header  *Header

	// Service to be invoked.
	Service string `xml:"service"`
}

// Ticket header.
type Header struct {
	XMLName        xml.Name `xml:"header"`
	UniqueID       uint32   `xml:"uniqueId"`
	GenerationTime string   `xml:"generationTime"`
	ExpirationTime string   `xml:"expirationTime"`
}

// CreateTicketRequest creates an xml request with generation, expiration, uniqueID
// and invoked service.
func CreateTicketRequest(service string, exp int16) ([]byte, error) {
	uid := rand.Uint32()

	// Dt gets the time synchronized with the AFIP/ARCA servers.
	dt, err := ntp.Time("time.afip.gov.ar")
	if err != nil {
		return nil, fmt.Errorf("error getting ntp time: %v", err)
	}

	// Create generation and expire time.
	gt := dt.Format("2006-01-02") + "T" + dt.Format("15:04:05")
	timeIn := dt.Add(time.Hour * time.Duration(exp))
	ed := timeIn.Format("2006-01-02") + "T" + timeIn.Format("15:04:05")

	// Create ticket request XML.
	newHead := &Header{UniqueID: uid, GenerationTime: gt, ExpirationTime: ed}
	newReq := &TicketRequest{Header: newHead, Service: service}

	// Format and return XML
	out, err := xml.MarshalIndent(newReq, " ", " ")
	if err != nil {
		return nil, err
	}

	return out, err
}
