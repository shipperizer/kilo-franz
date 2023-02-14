package tls

import (
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

func TestGetTLSWithSM(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	client, err := SMClient(specs.Region, specs.Endpoint)

	assert := assert.New(t)
	assert.Nil(err)

	smConfig := SecretManagerConfig{
		CertificateString: "test/cert",
		KeyString:         "test/key",
		SMClient:          client,
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestGetTLSWithSMP12Sha256(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	client, err := SMClient(specs.Region, specs.Endpoint)

	assert := assert.New(t)
	assert.Nil(err)

	smConfig := SecretManagerConfig{
		P12String: "test/p12.sha256",
		SMClient:  client,
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		UseP12:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestGetTLSWithSMP12(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	client, err := SMClient(specs.Region, specs.Endpoint)

	assert := assert.New(t)
	assert.Nil(err)

	smConfig := SecretManagerConfig{
		P12String: "test/p12",
		SMClient:  client,
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		UseP12:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestGetTLSWithSMFails(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	client, err := SMClient(specs.Region, specs.Endpoint)

	assert := assert.New(t)
	assert.Nil(err)

	smConfig := SecretManagerConfig{
		CertificateString: "fake/cert",
		KeyString:         "fake/key",
		SMClient:          client,
	}

	cfg := TLSConfig{
		UseTLS:                  true,
		SMConfig:                &smConfig,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	tls, err := GetTLS(cfg)

	assert.Nil(tls)
	assert.NotNil(err)
}

func TestGetTLSWithSMNoSmConfig(t *testing.T) {
	cfg := TLSConfig{
		UseTLS:                  true,
		SMConfig:                nil,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	assert := assert.New(t)

	_, err := GetTLS(cfg)

	assert.NotNil(err)
}

func TestGetTLSWithSMNoSmConfigUseTLSFalse(t *testing.T) {
	cfg := TLSConfig{
		UseTLS:                  false,
		SMConfig:                nil,
		ClientSignedCertificate: []byte(""),
		ClientKey:               []byte(""),
	}

	assert := assert.New(t)

	c, err := GetTLS(cfg)

	assert.Nil(c)
	assert.Nil(err)
}
