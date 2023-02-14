package vault

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// SecretsManagerAPI interface for AWS Secrets Manager Client.
type SecretsManagerAPI interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}
