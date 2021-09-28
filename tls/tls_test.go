package tls

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetTLSWithSM(t *testing.T) {
	smConfig := SecretManagerConfig{
		CertificateString: "test/cert",
		KeyString:         "test/key",
		SMClient:          SMTestClient(),
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert := assert.New(t)

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestGetTLSWithSMFails(t *testing.T) {
	smConfig := SecretManagerConfig{
		CertificateString: "fake/cert",
		KeyString:         "fake/key",
		SMClient:          SMTestClient(),
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert := assert.New(t)

	assert.Nil(tls)
	assert.NotNil(err)
}
