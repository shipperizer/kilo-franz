package sasl

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/plain"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/shipperizer/kilo-franz/logging"
)

// TODO @shipperizer allow to shape structure of secrets and default to this
type SecretCredentials struct {
	Keys    string `json:"api-keys"`
	Secrets string `json:"api-secrets"`
}

// SecretManagerConfig holds the key values to fetch from SM for the credentials key
type SecretManagerConfig struct {
	CredentialsKey string

	Vault VaultInterface
}

func NewSecretManagerConfig(credentialsKey string, vault VaultInterface) *SecretManagerConfig {
	c := new(SecretManagerConfig)

	c.CredentialsKey = credentialsKey
	c.Vault = vault

	return c
}

type SASLConfig struct {
	username  string
	password  string
	saslType  SASLType
	algorithm scram.Algorithm

	smConfig *SecretManagerConfig
	logger   logging.LoggerInterface
}

func (c *SASLConfig) getCredentials(ctx context.Context) (string, string) {
	if c.smConfig == nil {
		return c.username, c.password
	}

	credentials, err := c.smConfig.Vault.GetValue(ctx, c.smConfig.CredentialsKey)

	if err != nil {
		c.logger.Errorf("issues retrieving SASL credentials secret %s: %s", c.smConfig.CredentialsKey, err)
		c.logger.Warn("proceed with default credentials")

		return c.username, c.password
	}

	secrets := SecretCredentials{}

	err = json.Unmarshal(credentials, &secrets)

	if err != nil {
		c.logger.Errorf("issues deserializing SASL secret %s: %s", c.smConfig.CredentialsKey, err)
		c.logger.Warn("proceed with default credentials")

		return c.username, c.password
	}

	return secrets.Keys, secrets.Secrets
}

func (c *SASLConfig) GetSASLMechanism() sasl.Mechanism {
	username, password := c.getCredentials(context.TODO())

	switch c.saslType {
	case PlainSASL:
		return plain.Mechanism{
			Username: username,
			Password: password,
		}
	case ScramSASL:
		mechanism, err := scram.Mechanism(c.algorithm, username, password)

		if err != nil {
			return nil
		}

		return mechanism
	default:
		c.logger.Errorf("unknown SASL type %s, defaulting to no SASL", c.saslType)
		return nil
	}
}

func NewSASLConfig(username, password string, saslType SASLType, algorithm scram.Algorithm, smConfig *SecretManagerConfig, logger logging.LoggerInterface) *SASLConfig {
	c := new(SASLConfig)

	c.username = username
	c.password = password
	c.saslType = saslType
	c.algorithm = algorithm
	c.smConfig = smConfig
	c.logger = logger

	return c
}
