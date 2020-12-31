package main

import (
	"net/http"
	"strconv"
	"time"

	"go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	app.GET("/datetime", dateTimeHandler)

	srv := &http.Server{
		Handler: app,
		Addr:    ":8082",
		TLSConfig: secure.MustLoadServerTLS(
			"ca.pem",
			"service-timestamps.pem",
			"service-timestamps-key.pem",
		),
	}

	srv.ListenAndServeTLS("", "")
}

func dateTimeHandler(c *gin.Context) {
	t, err := strconv.ParseInt(c.Query("timestamp"), 10, 64)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{"datetime": time.Unix(t, 0)})
}
