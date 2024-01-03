package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"sync"
)

// SafeCert is a struct to hold and manage a tls certificate
type SafeCert struct {
	keyPem  []byte
	certPem []byte

	cert *tls.Certificate

	sync.RWMutex
}

// newSafeCert makes a SafeCert using the supplied tlsCert
func NewSafeCert() *SafeCert {
	return &SafeCert{}
}

// TlsCertFunc returns the function to get the tls.Certificate from SafeCert
func (sc *SafeCert) TlsCertFunc() func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
	return func(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		sc.RLock()
		defer sc.RUnlock()

		return sc.cert, nil
	}
}

// Update updates the certificate with the specified key and cert pem
func (sc *SafeCert) Update(keyPem, certPem []byte) (updated bool, err error) {
	sc.Lock()
	defer sc.Unlock()

	// if no update to do, don't do anything
	if bytes.Equal(sc.keyPem, keyPem) && bytes.Equal(sc.certPem, certPem) {
		return false, nil
	}

	// update pem in cert struct
	sc.keyPem = keyPem
	sc.certPem = certPem

	// make tls certificate
	tlsCert, err := tls.X509KeyPair(certPem, keyPem)
	if err != nil {
		return false, fmt.Errorf("failed to make x509 key pair for cert update (%s)", err)
	}

	// update certificate
	sc.cert = &tlsCert

	return true, nil
}
