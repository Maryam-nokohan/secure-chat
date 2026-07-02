package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
	"io"
)

type EncryptionService struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

type EncryptedData struct {
	EncryptedKey []byte
	Nonce        []byte
	Ciphertext   []byte
}

func NewEncryptionService() (*EncryptionService, error) {
	privatKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		LogError(err)
		return nil, err
	}

	publickKey := &privatKey.PublicKey

	return &EncryptionService{
		privateKey: privatKey,
		publicKey:  publickKey,
	}, nil
}

func (e *EncryptionService) Encrypt(plain []byte) (*EncryptedData, error) {
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
		e.publicKey,
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

func (e *EncryptionService) Decrypt(data *EncryptedData) ([]byte, error) {
	aesKey, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		e.privateKey,
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

	plaintext, err := gcm.Open(
		nil,
		data.Nonce,
		data.Ciphertext,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("decryption failed: %w", err)
	}

	return plaintext, nil
}
