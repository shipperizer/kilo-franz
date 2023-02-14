package vault

import (
	"context"
	"testing"

	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func TestAWSVault(t *testing.T) {
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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	assert := assert.New(t)

	assert.Nil(err)

	client := secretsmanager.NewFromConfig(cfg)

	vault := NewAWSVault(client)

	assert.NotNil(vault)
}

func TestAWSVaultGetValueStringSuccessful(t *testing.T) {
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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	assert := assert.New(t)

	assert.Nil(err)

	client := secretsmanager.NewFromConfig(cfg)

	vault := NewAWSVault(client)

	password, err := vault.GetValue(context.Background(), "test/password")

	assert.Nil(err, "Error fetching password")
	assert.Equal([]byte("1234"), password)
}

func TestAWSVaultGetValueStringNotFound(t *testing.T) {
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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	assert := assert.New(t)

	assert.Nil(err)

	client := secretsmanager.NewFromConfig(cfg)

	vault := NewAWSVault(client)
	_, err = vault.GetValue(context.Background(), "fake/password")

	assert.NotNil(err, "Should have not found the secret")
}

func TestAWSVaultGetValueBinarySuccessful(t *testing.T) {
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

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithRegion(specs.Region),
	)

	assert := assert.New(t)

	assert.Nil(err)

	client := secretsmanager.NewFromConfig(cfg)

	vault := NewAWSVault(client)

	binary, err := vault.GetValue(context.Background(), "test/bin")

	assert.Nil(err, "Error fetching password")
	assert.Equal([]byte("1234"), binary)
}

func TestSMTestClientEndpointFail(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(specs.Region),
	)

	assert := assert.New(t)

	assert.Nil(err)

	client := secretsmanager.NewFromConfig(cfg)

	vault := NewAWSVault(client)

	_, err = vault.GetValue(context.Background(), "fake/password")

	assert.NotNil(err)
}
