package tls

import (
	"context"
	_tls "crypto/tls"
	"fmt"
	"time"

	"github.com/shipperizer/kilo-franz/logging"
	"golang.org/x/sync/errgroup"
)

type SecretManagerConfig struct {
	CertificateString string
	KeyString         string
	P12String         string
	SMClient          SecretsManagerAPI // TODO @shipperizer deprecate in favour of the one below
	Vault             VaultInterface
}

func NewSecretManagerConfig(cert, key, p12 string, vault VaultInterface) *SecretManagerConfig {
	c := new(SecretManagerConfig)

	c.CertificateString = cert
	c.KeyString = key
	c.P12String = p12
	c.Vault = vault

	return c
}

// TLSConfig holds core configuration to setup TLS for kafka
type TLSConfig struct {
	UseTLS                  bool
	UseP12                  bool
	SMConfig                *SecretManagerConfig
	ClientSignedCertificate []byte
	ClientKey               []byte

	logger logging.LoggerInterface
}

func (c *TLSConfig) GetTLS() (*_tls.Config, error) {
	if !c.UseTLS {
		return nil, nil
	}

	var cert, key []byte
	var err error

	if c.SMConfig == nil {
		cert = c.ClientSignedCertificate
		key = c.ClientKey

		c.logger.Debug("Client Signed Cert: ", string(cert))
		c.logger.Debug("Client Key: ", string(key))

		return MakeTLS(cert, key, c.UseP12)
	}

	if c.SMConfig.Vault == nil {
		return GetTLS(*c)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get keys concurrently
	g, ctx := errgroup.WithContext(ctx)

	if c.UseP12 {
		g.Go(func() error {
			var err error
			cert, err = c.SMConfig.Vault.GetValue(ctx, c.SMConfig.P12String)
			c.logger.Debugf("certificate: %s", cert)
			return err
		})
	} else {
		g.Go(func() error {
			var err error
			cert, err = c.SMConfig.Vault.GetValue(ctx, c.SMConfig.CertificateString)
			c.logger.Debugf("certificate: %s", cert)
			return err
		})

		g.Go(func() error {
			var err error
			key, err = c.SMConfig.Vault.GetValue(ctx, c.SMConfig.KeyString)
			c.logger.Debugf("key: %s", key)
			return err
		})
	}

	err = g.Wait() // wait and check for error

	if err != nil {
		return nil, fmt.Errorf("issues fetching SM values: %s", err)
	}

	return MakeTLS(cert, key, c.UseP12)

}

func NewTLSConfig(cert, key []byte, useTLS, useP12 bool, secretsConfig *SecretManagerConfig, logger logging.LoggerInterface) *TLSConfig {
	c := new(TLSConfig)

	c.UseTLS = useTLS
	c.UseP12 = useP12
	c.SMConfig = secretsConfig
	c.ClientSignedCertificate = cert
	c.ClientKey = key
	c.logger = logger

	return c
}
