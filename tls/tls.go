package tls

import (
	"context"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	pkcs12 "software.sslmate.com/src/go-pkcs12"
)

// GetTLS produces a TLS object to be used by kafka consumer/producer
func GetTLS(cfg TLSConfig) (*tls.Config, error) {
	if !cfg.UseTLS {
		return nil, nil
	}

	var cert, key []byte
	var err error

	if cfg.SMConfig == nil {
		cert = cfg.ClientSignedCertificate
		key = cfg.ClientKey

		log.Debug("Client Signed Cert: ", string(cert))
		log.Debug("Client Key: ", string(key))
	} else {
		sm := cfg.SMConfig.SMClient

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// get keys concurrently
		g, ctx := errgroup.WithContext(ctx)

		if cfg.UseP12 {
			g.Go(func() error {
				var err error
				cert, err = GetSMValue(ctx, sm, cfg.SMConfig.P12String)
				log.Debugf("Certificate: %s", cert)
				return err
			})
		} else {
			g.Go(func() error {
				var err error
				cert, err = GetSMValue(ctx, sm, cfg.SMConfig.CertificateString)
				log.Debugf("Certificate: %s", cert)
				return err
			})

			g.Go(func() error {
				var err error
				key, err = GetSMValue(ctx, sm, cfg.SMConfig.KeyString)
				log.Debugf("key: %s", key)
				return err
			})
		}

		err = g.Wait() // wait and check for error

		if err != nil {
			return nil, fmt.Errorf("issues fetching SM values: %s", err)
		}
	}

	var tls *tls.Config

	tls, err = MakeTLS(cert, key, cfg.UseP12)

	if err != nil {
		return nil, fmt.Errorf("issues with the TLS config: %s", err)
	}

	return tls, nil
}

// MakeTLS generates a tls.Config, kindly stolen from
// https://github.com/discovery-digital/entitlements-collection/blob/master/kafkaclient/client.go#L230
func MakeTLS(clientCert, key []byte, isP12 bool) (*tls.Config, error) {
	var err error

	// if isP12 is true then clientCert is in p12 format
	if isP12 {
		clientCert, key, err = DecodeP12(clientCert)

		if err != nil {
			return nil, err
		}
	}

	cert, err := tls.X509KeyPair(clientCert, key)

	if err != nil {
		return nil, err
	}

	log.Debugf("Key: %s", key)
	log.Debugf("Client Cert: %s", clientCert)
	log.Debugf("Cert: %v", cert)

	if err != nil {
		return nil, err
	}

	// Get the SystemCertPool, continue with an empty pool on error
	rootCAs, err := x509.SystemCertPool()

	if rootCAs == nil {
		rootCAs = x509.NewCertPool()
		log.Warnf("Using empty cert-pool because of - %v", err)
	} else {
		log.Info("Using system cert-pool")
	}

	for _, cert := range DecodePEM(clientCert).Certificate {
		x509Cert, err := x509.ParseCertificate(cert)
		if err != nil {
			log.Errorf("issue parsing cert PEM: %s", err.Error())
		}
		rootCAs.AddCert(x509Cert)
	}

	log.Debugf("RootCa: %v", rootCAs)
	log.Debugf("Certificates: %v", []tls.Certificate{cert})

	return &tls.Config{
		RootCAs:            rootCAs,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}, nil
}

// DecodePEM builds a PEM certificate object
func DecodePEM(certPEM []byte) tls.Certificate {
	var cert tls.Certificate
	var certDER *pem.Block
	for {
		certDER, certPEM = pem.Decode(certPEM)
		if certDER == nil {
			break
		}
		if certDER.Type == "CERTIFICATE" {
			cert.Certificate = append(cert.Certificate, certDER.Bytes)
		}
	}

	return cert
}

func DecodeP12(p12 []byte) ([]byte, []byte, error) {

	privateKey, cert, _, err := pkcs12.DecodeChain(p12, "")

	if err != nil {
		return nil, nil, err
	}

	rsaPk, ok := privateKey.(*rsa.PrivateKey)

	if !ok {
		return nil, nil, err
	}

	rsaPkBytes, err := x509.MarshalPKCS8PrivateKey(rsaPk)

	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: rsaPkBytes})
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})

	return certPEM, privateKeyPEM, nil
}
