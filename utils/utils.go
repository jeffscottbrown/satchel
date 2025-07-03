package utils

import (
	"log/slog"
	"os"

	"github.com/jeffscottbrown/gogoogle/secrets"
)

func RetrieveSecretValue(secretName string) string {
	clientSecret, err := secrets.RetrieveSecret(secretName)
	if err != nil {
		slog.Warn("Falling back to OS environment variable", slog.String("secretName", secretName), slog.Any("error", err))
		clientSecret = os.Getenv(secretName)
	}

	return clientSecret
}
