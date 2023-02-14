package sasl

import (
	"context"
)

// VaultInterface provides method(s) to fetch secrets
type VaultInterface interface {
	GetValue(ctx context.Context, key string) ([]byte, error)
}
