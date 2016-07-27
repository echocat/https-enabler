package main

import (
	"crypto/tls"
	"io"
	"log"
	"net/http"
)

func startHttpsServer(listenAddress, listenCert, listenPrivateKey, listenCa string) {
	proxy := &proxy{}

	targetListenCa := listenCa
	if len(targetListenCa) <= 0 {
		targetListenCa = listenCert
	}
	certificates, err := loadCertificatesFrom(targetListenCa)
	if err != nil {
		log.Fatalf("Couldn't load client CAs from %s. Got: %s", listenCa, err)
	}

	server := &http.Server{
		Addr:    listenAddress,
		Handler: proxy,
		TLSConfig: &tls.Config{
			ClientCAs:  certificates,
			ClientAuth: tls.RequireAndVerifyClientCert,
		},
	}

	targetPrivateKey := listenPrivateKey
	if len(targetPrivateKey) <= 0 {
		targetPrivateKey = listenCert
	}
	log.Printf("Listening on %s (scheme=HTTPS, secured=TLS, clientValidation=on)\n", server.Addr)
	err = server.ListenAndServeTLS(listenCert, targetPrivateKey)

	if err != nil {
		log.Fatalf("Could not start server. Got: %v", err)
	}
}

type proxy struct {
}

func (instance *proxy) ServeHTTP(writer http.ResponseWriter, in *http.Request) {
	client := &http.Client{}
	requestUrl := *in.URL
	requestUrl.Scheme = "http"
	requestUrl.Host = *connectAddress
	request := &http.Request{
		URL:    &requestUrl,
		Proto:  in.Proto,
		Method: in.Method,
		Header: in.Header,
		Host:   *connectAddress,
	}
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Could not handle request: %v %v. Cause: %v", in.Method, in.RequestURI, err)
		return
	}

	for key, values := range response.Header {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	io.Copy(writer, response.Body)
}
