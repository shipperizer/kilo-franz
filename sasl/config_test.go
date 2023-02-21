package sasl

import (
	"context"
	"testing"

	"github.com/kelseyhightower/envconfig"
	_sasl "github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/golang/mock/gomock"
	"github.com/shipperizer/kilo-franz/logging"
	"github.com/shipperizer/kilo-franz/vault"
)

//go:generate mockgen -build_flags=--mod=mod -package sasl -destination ./mock_logging.go -source=../logging/interfaces.go

func TestSASLConfigGetSASLMechanismPlain(t *testing.T) {
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

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := NewSASLConfig("", "", true, PlainSASL, nil, cfgSM, mockLogger)

	assert := assert.New(t)

	sasl := cfgSASL.GetSASLMechanism()

	assert.NotNil(sasl)
	assert.IsType(plain.Mechanism{}, sasl)
}

func TestSASLConfigGetSASLMechanismScram(t *testing.T) {
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

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := NewSASLConfig("", "", true, ScramSASL, scram.SHA256, cfgSM, mockLogger)

	assert := assert.New(t)

	sasl := cfgSASL.GetSASLMechanism()

	assert.NotNil(sasl)
	assert.Implements((*_sasl.Mechanism)(nil), sasl)
}

func TestSASLConfigGetCredentials(t *testing.T) {
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

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := NewSASLConfig("", "", true, ScramSASL, scram.SHA256, cfgSM, logging.NewLogger())

	assert := assert.New(t)

	username, password := cfgSASL.getCredentials(context.TODO())

	assert.Equal("admin", username)
	assert.Equal("admin-secret", password)
}

func TestSASLConfigGetCredentialsBadSecret(t *testing.T) {
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

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("sasl/fake", vaultAWS)
	cfgSASL := NewSASLConfig("default", "default", true, ScramSASL, scram.SHA256, cfgSM, mockLogger)

	assert := assert.New(t)

	username, password := cfgSASL.getCredentials(context.TODO())

	assert.Equal("default", username)
	assert.Equal("default", password)
}

func TestSASLConfigGetCredentialsNoSM(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMockLoggerInterface(ctrl)

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()

	cfgSASL := NewSASLConfig("default", "default", true, ScramSASL, scram.SHA256, nil, mockLogger)

	assert := assert.New(t)

	username, password := cfgSASL.getCredentials(context.TODO())

	assert.Equal("default", username)
	assert.Equal("default", password)
}

func TestSASLConfigGetSASLMechanismWithSASLDisabled(t *testing.T) {
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

	mockLogger.EXPECT().Errorf(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Warn(gomock.Any()).AnyTimes()
	mockLogger.EXPECT().Debug(gomock.Any()).AnyTimes()

	vaultAWS := vault.NewAWSVault(client)
	cfgSM := NewSecretManagerConfig("sasl/credentials", vaultAWS)
	cfgSASL := NewSASLConfig("", "", false, PlainSASL, nil, cfgSM, mockLogger)

	assert := assert.New(t)

	sasl := cfgSASL.GetSASLMechanism()

	assert.Nil(sasl)
}
