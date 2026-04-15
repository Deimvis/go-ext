package xcrypto

import (
	"crypto/rand"
	"io"
)

func GenerateSecretKey(size int) []byte {
	bytes := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, bytes)
	if err != nil {
		panic(err)
	}
	return bytes
}
