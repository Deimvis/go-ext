package xcrypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
)

// GenerateRSAKeys returns (private, public) keys.
// These keys are also called secret and public keys (sk and pk)
func GenerateRSAKeys() (secretKey []byte, publicKey []byte) {
	sk, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}
	pk := &sk.PublicKey

	skBytes := x509.MarshalPKCS1PrivateKey(sk)
	pkBytes := x509.MarshalPKCS1PublicKey(pk)
	return skBytes, pkBytes
}

// SignRSA return (signature, hash, error)
func SignRSA(skBytes []byte, msg []byte) ([]byte, []byte, error) {
	hash, err := HashSHA256(msg)
	if err != nil {
		return nil, nil, err
	}
	sk, err := x509.ParsePKCS1PrivateKey(skBytes)
	if err != nil {
		return nil, nil, err
	}
	sig, err := rsa.SignPKCS1v15(rand.Reader, sk, crypto.SHA256, hash)
	return sig, hash, err
}

func VerifyRSA(pkBytes []byte, sig []byte, hash []byte) error {
	pk, err := x509.ParsePKCS1PublicKey(pkBytes)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pk, crypto.SHA256, hash, sig)
}

func HashSHA256(msg []byte) ([]byte, error) {
	msgHash := sha256.New()
	_, err := msgHash.Write(msg)
	if err != nil {
		return nil, err
	}
	return msgHash.Sum(nil), nil
}
