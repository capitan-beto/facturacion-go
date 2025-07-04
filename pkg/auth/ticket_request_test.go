package auth

import (
	"testing"
)

func TestCreateTicketRequest(t *testing.T) {
	_, err := CreateTicketRequest("wsfe", 18)
	if err != nil {
		t.Fatal("error en la creación de la petición")
	}

}
