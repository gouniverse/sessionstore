package sessionstore

import (
	"github.com/gouniverse/strutils"
)

func generateSessionKey(keyLength int) string {
	gamma := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	return strutils.RandomFromGamma(keyLength, gamma)
}
