package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"go-mutual-tls/secure"

	"github.com/gin-gonic/gin"
)

var testClient *http.Client

func TestRunApp(t *testing.T) {
	run, shutdown := App(context.Background(), 8080, ioutil.Discard)
	defer shutdown()
	go run()

	r, err := testClient.Get("https://localhost:8080/datetime?timestamp=1609417705")
	if err != nil {
		t.Fatal(err)
	}

	defer r.Body.Close()
	body, err := ioutil.ReadAll(r.Body)

	var dr datetimeResponse
	if err := json.Unmarshal(body, &dr); err != nil {
		t.Fatal(err)
	}

	expected := "2020-12-31 14:28:25 +0200 EET"
	if dr.Datetime.String() != expected {
		t.Errorf("expected %s got %s", expected, dr.Datetime.String())
	}
}

type datetimeResponse struct {
	Datetime time.Time `json:"datetime"`
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	testClient = buildTestClient()
	code := m.Run()
	os.Exit(code)
}

func buildTestClient() *http.Client {
	return secure.NewHTTPSClient(
		secure.MustLoadClientTLS("ca.pem", "service-apigw.pem", "service-apigw-key.pem"),
	)
}
