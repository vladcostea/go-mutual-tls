package main

import (
	"net/http"

	"github.com/vladcostea/go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/version", func(c *gin.Context) {
		c.JSON(200, gin.H{"version": "0.0.1"})
	})

	srv := &http.Server{
		Handler:   r,
		Addr:      ":8083",
		TLSConfig: secure.MustLoadServerTLS("ca.pem", "service-version.pem", "service-version-key.pem"),
	}

	srv.ListenAndServeTLS("", "")
}
