package main

import (
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"os"
	"fmt"
)

func VerifyCert(dns_name string, ca_file string, server_cert_file string) {
	ca, err := ioutil.ReadFile(ca_file)
	if err != nil {
			panic("failed to read CA file: "  + err.Error())
	}

	server_cert, err := ioutil.ReadFile(server_cert_file)
	if err != nil {
			panic("failed to read server cert file: "  + err.Error())
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(ca))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(server_cert))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		DNSName: dns_name,
		Roots:   roots,
	}

	if _, err := cert.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}

	fmt.Println("OK")
}

func main() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: " + os.Args[0] + " dns-name-to-verify.example.org ca.crt server_cert.crt")
		os.Exit(1)
	}
	VerifyCert(os.Args[1], os.Args[2], os.Args[3])
}
