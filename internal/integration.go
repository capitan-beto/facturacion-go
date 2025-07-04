package internal

import (
	"cmd/api/main.go/pkg/auth"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func IntegrationAuth() {
	// Step 1: Create auth xml (TRA)

	at, err := auth.CreateTicketRequest("wsfe", 18)
	if err != nil {
		log.Error(err)
		return
	}

	//Step 2: Create CMS
	loginCMS, err := auth.GenerateCMS(at, "MiCertificado.pem", "MiClavePrivada.key")
	if err != nil {
		log.Error(err)
		return
	}

	// Step 3: Create login request
	loginRequest, err := auth.CreateLoginRequest(loginCMS)
	if err != nil {
		log.Error(err)
		return
	}

	// Step 4: Make request to AFIP WSAA.

	res, err := auth.RequestAuth(loginRequest, "https://wsaahomo.afip.gov.ar/ws/services/LoginCms")
	if err != nil {
		log.Error(err)
		return
	}

	// Step 5: Parse response

	authMsg, err := auth.ParseResponse(res)
	if err != nil {
		errMsg, err := auth.ParseResponseError(res)
		if err != nil {
			log.Error(err)
			return
		}
		fmt.Println(errMsg)
	} else {
		fmt.Println(authMsg)
	}

}
