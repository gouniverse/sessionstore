package sessionstore

import (
	"os"

	"github.com/gouniverse/strutils"
)

func generateSessionKey(keyLength int) string {
	gamma := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return strutils.RandomFromGamma(keyLength, gamma)
}

// fileExists checks if a file exists
func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)

	return !os.IsNotExist(err)
}
