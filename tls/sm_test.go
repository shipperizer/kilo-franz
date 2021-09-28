package tls

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func SMTestClient() SecretsManagerAPI {
	viper.BindEnv("aws.region", "AWS_REGION")
	viper.BindEnv("aws.endpoint", "AWS_ENDPOINT")

	smApi, _ := SMClient(viper.GetString("aws.region"), viper.GetString("aws.endpoint"))

	return smApi
}

func TestSMClient(t *testing.T) {
	viper.BindEnv("aws.region", "AWS_REGION")
	viper.BindEnv("aws.endpoint", "AWS_ENDPOINT")

	client, err := SMClient(viper.GetString("aws.region"), viper.GetString("aws.endpoint"))

	assertion := assert.New(t)

	var smApi *secretsmanager.Client

	assertion.Nil(err)
	assertion.NotNil(client)
	assertion.IsType(smApi, client, "Should have been SecretsManagerAPI interface")
}

func TestGetSMValueStringSuccessful(t *testing.T) {
	password, err := GetSMValue(context.Background(), SMTestClient(), "test/password")

	assert := assert.New(t)

	assert.Nil(err, "Error fetching password")
	assert.Equal([]byte("1234"), password)
}

func TestGetSMValueStringNotFound(t *testing.T) {
	_, err := GetSMValue(context.Background(), SMTestClient(), "fake/password")

	assertion := assert.New(t)
	assertion.NotNil(err, "Should have not found the secret")
}

func TestGetSMValueBinarySuccessful(t *testing.T) {
	binary, err := GetSMValue(context.Background(), SMTestClient(), "test/bin")

	assert := assert.New(t)

	assert.Nil(err, "Error fetching password")
	assert.Equal([]byte("1234"), binary)
}
