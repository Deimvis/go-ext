package xcrypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

func EncryptAES256CBC(text string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := []byte(text)
	b = PKCS7Padding(b, aes.BlockSize)
	encrypted := make([]byte, aes.BlockSize+len(b))
	iv := encrypted[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(encrypted[aes.BlockSize:], b)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptAES256CBC(encryptedText string, key []byte) (string, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	if len(encryptedData) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	iv := encryptedData[:aes.BlockSize]
	encryptedData = encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(encryptedData, encryptedData)

	decryptedData, err := PKCS7Unpadding(encryptedData, aes.BlockSize)
	if err != nil {
		return "", err
	}

	return string(decryptedData), nil
}

func PKCS7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS7Unpadding(paddedtext []byte, blockSize int) ([]byte, error) {
	length := len(paddedtext)
	if length == 0 {
		return nil, fmt.Errorf("invalid padding size")
	}

	padding := int(paddedtext[length-1])
	if padding > blockSize || padding > length {
		return nil, fmt.Errorf("invalid padding")
	}

	for _, v := range paddedtext[length-padding:] {
		if int(v) != padding {
			return nil, fmt.Errorf("invalid padding")
		}
	}

	return paddedtext[:length-padding], nil
}

func GenerateAES256SecretKey() []byte {
	return GenerateSecretKey(32)
}
