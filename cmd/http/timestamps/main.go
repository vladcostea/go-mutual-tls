package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

func main() {
	run, shutdown := App(context.Background(), 8082, os.Stdout)
	defer shutdown()
	run()
}

func App(ctx context.Context, port int, stdout io.Writer) (func() error, func() error) {
	app := gin.New()
	app.Use(gin.Recovery())
	app.Use(gin.LoggerWithWriter(stdout))
	app.GET("/datetime", dateTimeHandler)

	srv := &http.Server{
		Handler: app,
		Addr:    fmt.Sprintf(":%d", port),
		TLSConfig: secure.MustLoadServerTLS(
			"ca.pem",
			"service-timestamps.pem",
			"service-timestamps-key.pem",
		),
	}

	return func() error {
			return srv.ListenAndServeTLS("", "")
		}, func() error {
			return srv.Shutdown(ctx)
		}
}

func dateTimeHandler(c *gin.Context) {
	t, err := strconv.ParseInt(c.Query("timestamp"), 10, 64)
	if err != nil {
		panic(err)
	}

	c.JSON(200, gin.H{"datetime": time.Unix(t, 0)})
}
