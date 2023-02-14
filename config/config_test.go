package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	// @shipperizer lesson learnt: don't name packages so generic, mocking won't like it
	_tls "github.com/shipperizer/kilo-franz/tls"
)

func TestNewConfigImplementsInterface(t *testing.T) {
	cfg := NewConfig(1*time.Hour, nil, nil, nil)

	assert := assert.New(t)

	assert.Implements((*ConfigInterface)(nil), cfg)
}

func TestNewConfigImplementsInterfaceWithTlsConfigFail(t *testing.T) {
	smClient, err := _tls.SMClient("region", "endpoint")

	if err != nil {
		t.Fatal(err)
	}

	tlsCfg := &_tls.TLSConfig{
		UseTLS: true,
		SMConfig: &_tls.SecretManagerConfig{
			CertificateString: "cert",
			KeyString:         "key",
			SMClient:          smClient,
		},
		ClientSignedCertificate: []byte(`cert`),
		ClientKey:               []byte(`key`),
	}

	cfg := NewConfig(1*time.Hour, tlsCfg, nil, nil)

	assert := assert.New(t)

	assert.NotNil(cfg)
	assert.Nil(cfg.getTLS(), "TLS with bad cert")
}
