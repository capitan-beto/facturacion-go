package auth

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func RequestAuth(b []byte, envURL string) ([]byte, error) {
	r := bytes.NewReader(b)

	req, err := http.NewRequest("POST", envURL, r)
	if err != nil {
		return nil, fmt.Errorf("error creating auth http request: %v", err)
	}

	req.Header.Add("Content-Type", "text/html;charset=UTF-8")
	req.Header.Add("SOAPAction", "urn:LoginCms")
	req.Header.Add("User-Agent", "Apache-HttpClient/4.1.1 (java 1.5)")

	transport := &http.Transport{
		Proxy: nil,
	}

	client := &http.Client{
		Transport: transport,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending auth http request: %v", err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading http body response: %v", err)
	}

	return body, err
}
