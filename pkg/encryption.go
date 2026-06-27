package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
)

type HybridCryptoService struct{}

type EncryptedPayload struct {
	Ciphertext   string
	EncryptedKey string
	IV           string
}

func NewHybridCryptoService() *HybridCryptoService {
	return &HybridCryptoService{}
}

func (s *HybridCryptoService) Encrypt(publicKeyPEM string, plaintext []byte) (*EncryptedPayload, error) {
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the public key")
	}
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}

	aesKey := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, aesKey); err != nil {
		return nil, err
	}

	blockCipher, err := aes.NewCipher(aesKey)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}
	iv := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	ciphertext := gcm.Seal(nil, iv, plaintext, nil)

	encryptedAESKey, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, aesKey, nil)
	if err != nil {
		return nil, err
	}

	return &EncryptedPayload{
		Ciphertext:   base64.StdEncoding.EncodeToString(ciphertext),
		EncryptedKey: base64.StdEncoding.EncodeToString(encryptedAESKey),
		IV:           base64.StdEncoding.EncodeToString(iv),
	}, nil
}

func (s *HybridCryptoService) Decrypt(privateKeyPEM string, payload *EncryptedPayload) ([]byte, error) {
    block, _ := pem.Decode([]byte(privateKeyPEM))
    if block == nil {
        return nil, errors.New("failed to parse PEM block")
    }
    privKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
    if err != nil {
        return nil, err
    }

    encryptedKeyBytes, err := base64.StdEncoding.DecodeString(payload.EncryptedKey)
    if err != nil {
        return nil, fmt.Errorf("decode encrypted key: %w", err)
    }
    ciphertextBytes, err := base64.StdEncoding.DecodeString(payload.Ciphertext)
    if err != nil {
        return nil, fmt.Errorf("decode ciphertext: %w", err)
    }
    ivBytes, err := base64.StdEncoding.DecodeString(payload.IV)
    if err != nil {
        return nil, fmt.Errorf("decode IV: %w", err)
    }

    aesKey, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, encryptedKeyBytes, nil)
    if err != nil {
        return nil, err
    }
    blockCipher, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(blockCipher)
    if err != nil {
        return nil, err
    }
    return gcm.Open(nil, ivBytes, ciphertextBytes, nil)
}
func ValidateRSAPublicKey(pemStr string) error {
    block, _ := pem.Decode([]byte(pemStr))
    if block == nil {
        return errors.New("invalid public key: not a PEM block")
    }
    pub, err := x509.ParsePKIXPublicKey(block.Bytes)
    if err != nil {
        return errors.New("invalid public key: " + err.Error())
    }
    if _, ok := pub.(*rsa.PublicKey); !ok {
        return errors.New("invalid public key: not an RSA key")
    }
    return nil
}