package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

type EncryptionService struct{}

type EncryptedData struct {
	EncryptedKey []byte
	Nonce        []byte
	Ciphertext   []byte
}

func NewEncryptionService() *EncryptionService {
	return &EncryptionService{}
}

func (EncryptionService) Encrypt(
	plain []byte,
	publicKey *rsa.PublicKey,
) (*EncryptedData, error) {

	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nil, nonce, plain, nil)

	encryptedKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		publicKey,
		aesKey,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &EncryptedData{
		EncryptedKey: encryptedKey,
		Nonce:        nonce,
		Ciphertext:   ciphertext,
	}, nil
}

func (EncryptionService) Decrypt(
	data *EncryptedData,
	privateKey *rsa.PrivateKey,
) ([]byte, error) {

	aesKey, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privateKey,
		data.EncryptedKey,
		nil,
	)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plain, err := gcm.Open(
		nil,
		data.Nonce,
		data.Ciphertext,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plain, nil
}

func GenerateRSAKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}
func EncodePublicKey(key *rsa.PublicKey) (string, error) {

	bytes, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		return "", err
	}

	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bytes,
	}

	return string(pem.EncodeToMemory(block)), nil
}

func EncodePrivateKey(key *rsa.PrivateKey) string {

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	return string(pem.EncodeToMemory(block))
}

func ParsePublicKey(pemString string) (*rsa.PublicKey, error) {

	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("invalid public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	return rsaPub, nil
}

func ParsePrivateKey(pemString string) (*rsa.PrivateKey, error) {

	block, _ := pem.Decode([]byte(pemString))
	if block == nil {
		return nil, errors.New("invalid private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}