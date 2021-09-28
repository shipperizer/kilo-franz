package tls

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	log "github.com/sirupsen/logrus"
)

// SecretsManagerAPI interface for AWS Secrets Manager Client.
type SecretsManagerAPI interface {
	GetSecretValue(ctx context.Context,
		params *secretsmanager.GetSecretValueInput,
		optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SMClient returns SecretsManagerAPI interface that implements secretsmanager.Client.
// `endpoint` arg can be used when developing locally.
func SMClient(region, endpoint string) (SecretsManagerAPI, error) {

	customResolver := aws.EndpointResolverFunc(func(service, r string) (aws.Endpoint, error) {
		if service == secretsmanager.ServiceID && endpoint != "" {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}
		// Returning EndpointNotFoundError will allow the service to fallback to it's default resolution.
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	cfg, err := awsconfig.LoadDefaultConfig(
		context.TODO(),
		awsconfig.WithEndpointResolver(customResolver),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return secretsmanager.NewFromConfig(cfg), nil
}

func GetSMValue(ctx context.Context, secretsManager SecretsManagerAPI, key string) ([]byte, error) {
	secret, err := secretsManager.GetSecretValue(
		ctx,
		&secretsmanager.GetSecretValueInput{
			SecretId: aws.String(key),
		},
	)

	if err != nil {
		return nil, err
	}

	if secret.SecretString != nil {
		return []byte(*secret.SecretString), nil
	}

	if secret.SecretBinary != nil {
		decodedSecret := make([]byte, base64.StdEncoding.DecodedLen(len(secret.SecretBinary)))
		len, err := base64.StdEncoding.Decode(decodedSecret, secret.SecretBinary)
		if err != nil {
			log.Error("Base64 Decode Error:", err)
			return nil, err
		}
		return decodedSecret[:len], nil
	}

	return nil, fmt.Errorf("Unexpected error fetching secret %s", key)
}
