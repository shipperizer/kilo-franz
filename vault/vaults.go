package vault

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type AWSVault struct {
	client SecretsManagerAPI
}

func (v *AWSVault) GetValue(ctx context.Context, key string) ([]byte, error) {
	secret, err := v.client.GetSecretValue(
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

func NewAWSVault(client SecretsManagerAPI) *AWSVault {
	v := new(AWSVault)
	v.client = client

	return v
}
