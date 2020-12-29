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
	cp.AppendCertsFromPEM(mustReadTLSFile(caFile))

	certs, err := tls.LoadX509KeyPair(tlsFilepath(certFile), tlsFilepath(keyFile))
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
	cp.AppendCertsFromPEM(mustReadTLSFile(caFile))

	certs, err := tls.LoadX509KeyPair(tlsFilepath(certFile), tlsFilepath(keyFile))
	if err != nil {
		panic(err)
	}

	return &tls.Config{
		Certificates: []tls.Certificate{certs},
		RootCAs:      cp,
	}
}

func mustReadTLSFile(file string) []byte {
	b, err := ioutil.ReadFile(tlsFilepath(file))
	if err != nil {
		panic(err)
	}

	return b
}

func tlsFilepath(file string) string {
	return filepath.Join(os.Getenv("TLS_DIR"), file)
}
