package tls

// SecretManagerConfig holds the key values to fetch from SM for the client cert and key
type SecretManagerConfig struct {
	CertificateString string
	KeyString         string
	SMClient          SecretsManagerAPI
}

// TLSConfig holds core configuration to setup TLS for kafka
type TLSConfig struct {
	UseTLS                  bool
	SMConfig                *SecretManagerConfig
	ClientSignedCertificate []byte
	ClientKey               []byte
}
