package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

func main() {
	client := secure.NewHTTPSClient(
		secure.MustLoadClientTLS("ca.pem", "service-apigw.pem", "service-apigw-key.pem"),
	)

	app := gin.Default()
	app.GET("/datetime", dateTimeHandler(client))

	srv := &http.Server{
		Handler: app,
		Addr:    ":8081",
		TLSConfig: secure.MustLoadServerTLS(
			"ca.pem",
			"service-apigw.pem",
			"service-apigw-key.pem",
		),
	}

	srv.ListenAndServeTLS("", "")
}

func dateTimeHandler(client *http.Client) gin.HandlerFunc {
	host := timestampsHost()
	return func(c *gin.Context) {
		url := fmt.Sprintf("%s/datetime?timestamp=%s", host, c.Query("timestamp"))
		r, err := client.Get(url)
		if err != nil {
			handle500(c, err)
			return
		}

		defer r.Body.Close()
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			handle500(c, err)
			return
		}

		var dr datetimeResponse
		err = json.Unmarshal(body, &dr)
		if err != nil {
			handle500(c, err)
			return
		}

		c.JSON(200, dr)
	}
}

func handle500(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
}

type datetimeResponse struct {
	Datetime time.Time `json:"datetime"`
}

func timestampsHost() string {
	host := os.Getenv("SERVICE_TIMESTAMPS_HOST")
	if host != "" {
		return host
	}

	return "https://localhost:8082"
}
