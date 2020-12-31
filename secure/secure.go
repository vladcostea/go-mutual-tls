package secure

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"os"
	"path/filepath"
)

func MustLoadServerTLS(caFile, certFile, keyFile string) *tls.Config {
	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM(mustReadFile(caFile))

	certs, err := tls.LoadX509KeyPair(certFilepath(certFile), certFilepath(keyFile))
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certs},
		ClientCAs:    cp,
		ClientAuth:   tls.RequireAndVerifyClientCert,
	}
}

func MustLoadClientTLS(caFile, certFile, keyFile string) *tls.Config {
	cp := x509.NewCertPool()
	cp.AppendCertsFromPEM(mustReadFile(caFile))

	certs, err := tls.LoadX509KeyPair(certFilepath(certFile), certFilepath(keyFile))
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certs},
		RootCAs:      cp,
	}
}

func mustReadFile(file string) []byte {
	b, err := ioutil.ReadFile(certFilepath(file))
	if err != nil {
		panic(err)
	}

	return b
}

func certFilepath(file string) string {
	return filepath.Join(os.Getenv("CERT_DIR"), file)
}
