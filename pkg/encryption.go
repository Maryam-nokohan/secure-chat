package pkg

import (

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	"errors"
)

type EncryptionService struct {
}

func NewEncryptionService() *EncryptionService {
    return &EncryptionService{}
}

func (s *EncryptionService) EncryptMessage(publicKeyPEM, plaintext string) (string, error) {
    pubKey, err := x509.ParsePKIXPublicKey([]byte(publicKeyPEM))
    if err != nil {
        return "", err
    }

    rsaPubKey, ok := pubKey.(*rsa.PublicKey)
    if !ok {
        return "", errors.New("invalid public key")
    }

    ciphertext, err := rsa.EncryptOAEP(
        sha256.New(),
        rand.Reader,
        rsaPubKey,
        []byte(plaintext),
        nil,
    )
    if err != nil {
        return "", err
    }

    return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *EncryptionService) DecryptMessage(privateKeyPEM, ciphertext string) (string, error) {
    privKey, err := x509.ParsePKCS1PrivateKey([]byte(privateKeyPEM))
    if err != nil {
        return "", err
    }

    ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
    if err != nil {
        return "", err
    }

    plaintext, err := rsa.DecryptOAEP(
        sha256.New(),
        rand.Reader,
        privKey,
        ciphertextBytes,
        nil,
    )
    if err != nil {
        return "", err
    }

    return string(plaintext), nil
}