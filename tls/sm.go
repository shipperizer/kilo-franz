package tls

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerAPI interface for AWS Secrets Manager Client.
type SecretsManagerAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

// SMClient returns SecretsManagerAPI interface that implements secretsmanager.Client.
// `endpoint` arg can be used when developing locally.
func SMClient(region, endpoint string) (SecretsManagerAPI, error) {

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, r string, opts ...interface{}) (aws.Endpoint, error) {
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
		awsconfig.WithEndpointResolverWithOptions(customResolver),
		awsconfig.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	return secretsmanager.NewFromConfig(cfg), nil
}

// TODO @shipperizer mvoe this to be a method and enhance the SecretsManagerAPI or split it and wrap it
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
		return secret.SecretBinary, nil
	}

	return nil, fmt.Errorf("unexpected error fetching secret %s", key)
}
