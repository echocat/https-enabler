package main

import (
	"crypto/x509"
	"io/ioutil"
)

func loadCertificatesFrom(pemFile string) (*x509.CertPool, error) {
	caCert, err := ioutil.ReadFile(pemFile)
	if err != nil {
		return nil, err
	}
	certificates := x509.NewCertPool()
	certificates.AppendCertsFromPEM(caCert)
	return certificates, nil
}
