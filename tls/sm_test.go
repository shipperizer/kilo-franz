package tls

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/kelseyhightower/envconfig"
	"github.com/stretchr/testify/assert"
)

func TestSMClient(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)

	client, err := SMClient(specs.Region, specs.Endpoint)

	assert := assert.New(t)

	var smApi *secretsmanager.Client

	assert.Nil(err)
	assert.NotNil(client)
	assert.IsType(smApi, client, "Should have been SecretsManagerAPI interface")
}

func TestGetSMValueStringSuccessful(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)
	assert := assert.New(t)

	client, err := SMClient(specs.Region, specs.Endpoint)
	assert.Nil(err)

	password, err := GetSMValue(context.Background(), client, "test/password")

	assert.Nil(err, "Error fetching password")
	assert.Equal([]byte("1234"), password)
}

func TestGetSMValueStringNotFound(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)
	assert := assert.New(t)

	client, err := SMClient(specs.Region, specs.Endpoint)
	assert.Nil(err)

	_, err = GetSMValue(context.Background(), client, "fake/password")

	assert.NotNil(err, "Should have not found the secret")
}

func TestGetSMValueBinarySuccessful(t *testing.T) {
	type EnvSpec struct {
		Region   string `envconfig:"aws_region"`
		Endpoint string `envconfig:"aws_sm_endpoint"`
	}

	var specs EnvSpec
	_ = envconfig.Process("", &specs)
	assert := assert.New(t)

	client, err := SMClient(specs.Region, specs.Endpoint)
	assert.Nil(err)

	binary, err := GetSMValue(context.Background(), client, "test/bin")

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
	assert := assert.New(t)

	client, err := SMClient(specs.Region, "")

	assert.Nil(err)

	_, err = GetSMValue(context.Background(), client, "fake/password")

	assert.NotNil(err)
}
