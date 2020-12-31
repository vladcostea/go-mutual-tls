package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/vladcostea/go-mutual-tls/secure"

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
		Addr:    ":8043",
		TLSConfig: secure.MustLoadServerTLS(
			"ca.pem",
			"service-apigw.pem",
			"service-apigw-key.pem",
		),
	}

	srv.ListenAndServeTLS("", "")
}

func dateTimeHandler(client *http.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		url := fmt.Sprintf("https://localhost:8083/datetime?timestamp=%s", c.Query("timestamp"))
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
