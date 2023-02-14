package tls

import (
	"context"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/golang/mock/gomock"
	"github.com/shipperizer/kilo-franz/vault"
)

//go:generate mockgen -build_flags=--mod=mod -package tls -destination ./mock_tls.go -source=interfaces.go
//go:generate mockgen -build_flags=--mod=mod -package tls -destination ./mock_logging.go -source=../logging/interfaces.go

func TestTLSConfigWithVault(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, options ...interface{}) (aws.Endpoint, error) {
		if service == secretsmanager.ServiceID && specs.Endpoint != "" {
			return aws.Endpoint{
				URL:           specs.Endpoint,
				SigningRegion: specs.Region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to it's default resolution.
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	client := secretsmanager.NewFromConfig(cfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("test/cert", "test/key", "", vaultAWS)
	cfgTLS := NewTLSConfig([]byte(""), []byte(""), true, false, cfgSM, mockLogger)

	assert := assert.New(t)

	tls, err := cfgTLS.GetTLS()

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestTLSConfigWithoutVault(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, options ...interface{}) (aws.Endpoint, error) {
		if service == secretsmanager.ServiceID && specs.Endpoint != "" {
			return aws.Endpoint{
				URL:           specs.Endpoint,
				SigningRegion: specs.Region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to it's default resolution.
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	client := secretsmanager.NewFromConfig(cfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)
	// logger is not called when not using the vault
	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).Times(0)

	cfgSM := SecretManagerConfig{
		CertificateString: "test/cert",
		KeyString:         "test/key",
		SMClient:          client,
	}

	cfgTLS := NewTLSConfig([]byte(""), []byte(""), true, false, &cfgSM, mockLogger)

	assert := assert.New(t)

	tls, err := cfgTLS.GetTLS()

	assert.NotNil(tls)
	assert.Nil(err)
}

func TestTLSConfigWithVaultP12(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, options ...interface{}) (aws.Endpoint, error) {
		if service == secretsmanager.ServiceID && specs.Endpoint != "" {
			return aws.Endpoint{
				URL:           specs.Endpoint,
				SigningRegion: specs.Region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to it's default resolution.
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, _ := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	client := secretsmanager.NewFromConfig(cfg)

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Debugf(gomock.Any(), gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("", "", "test/p12", vaultAWS)
	cfgTLS := NewTLSConfig([]byte(""), []byte(""), true, true, cfgSM, mockLogger)

	assert := assert.New(t)

	tls, err := cfgTLS.GetTLS()

	assert.NotNil(tls)
	assert.Nil(err)
}
