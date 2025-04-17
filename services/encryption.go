package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/scrypt"
)

// GenerateSalt generates a random salt of the specified length.
func GenerateSalt(length int) (string, error) {
	saltBytes := make([]byte, length)
	_, err := io.ReadFull(rand.Reader, saltBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate salt: %w", err)
	}
	salt := base64.StdEncoding.EncodeToString(saltBytes)
	return salt, nil
}

// DeriveKey derives a suitable encryption key from the user's PIN and salt using scrypt.
func DeriveKey(pin, salt string) ([]byte, error) {
	N := 16384   // CPU cost parameter
	r := 8       // Memory cost parameter
	p := 1       // Parallelization parameter
	keyLen := 32 // Desired key length (for AES-256)

	// Use the salt.
	key, err := scrypt.Key([]byte(pin+salt), []byte(salt), N, r, p, keyLen)
	if err != nil {
		return nil, fmt.Errorf("failed to derive key: %w", err)
	}
	return key, nil
}

// Encrypt encrypts the plaintext using the provided PIN and salt.
func Encrypt(pin, plaintext string) (string, error) {
	// Derive the encryption key from the PIN and salt.
	key, err := DeriveKey(pin, "") // Pass empty salt, the salt is already in the pin.
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the ciphertext using the provided PIN and salt.
func Decrypt(pin, ciphertext string) (string, error) {
	// Derive the decryption key from the PIN and salt.
	key, err := DeriveKey(pin, "") // Pass empty salt, the salt is already in the pin.
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	decodedCiphertext, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	if len(decodedCiphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertextBytes := decodedCiphertext[:nonceSize], decodedCiphertext[nonceSize:]

	plaintextBytes, err := aesGCM.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintextBytes), nil
}
