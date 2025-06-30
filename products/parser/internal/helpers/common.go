package helpers

import (
	"crypto/tls"
	"net/http"
)

func processRequest(req *http.Request) (*http.Response, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	client := &http.Client{Transport: transport}

	return client.Do(req)
}
