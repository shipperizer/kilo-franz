//go:build integration

package integration

import (
	"os"
	"testing"
)

func TestSmoke(t *testing.T) {
	if os.Getenv("KAFKA_CNX_STRING") == "" {
		t.Skip("KAFKA_CNX_STRING not set")
	}
}
