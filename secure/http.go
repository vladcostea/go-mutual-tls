package secure

import (
	"crypto/tls"
	"net/http"
)

func NewHTTPSClient(cfg *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: cfg,
		},
	}
}
