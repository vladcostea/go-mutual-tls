package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/vladcostea/go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

func main() {
	tlscfg := secure.MustLoadClientTLS("ca.pem", "client-apigw.pem", "client-apigw-key.pem")
	client := secure.NewHTTPSClient(tlscfg)

	app := gin.Default()
	app.GET("/version", versionHandler(client))

	srv := &http.Server{
		Handler:   app,
		Addr:      ":8043",
		TLSConfig: secure.MustLoadServerTLS("ca.pem", "service-apigw.pem", "service-apigw-key.pem"),
	}

	srv.ListenAndServeTLS("", "")
}

func versionHandler(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		r, err := client.Get("https://localhost:8083/version")
		return500IfError(c, err)

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		return500IfError(c, err)

		var vr versionResponse
		err = json.Unmarshal(body, &vr)
		return500IfError(c, err)

		c.JSON(200, vr)
	}
}

func return500IfError(c *gin.Context, err error) {
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
}

type versionResponse struct {
	Version string `json:"version"`
}
