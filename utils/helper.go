package utils

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateID() string {
	arr := make([]byte, 6)
	rand.Read(arr)
	return hex.EncodeToString(arr)
}
