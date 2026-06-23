// Package securecookie seals and unseals opaque cookie payloads with
// authenticated encryption (AES-256-GCM). It is used to carry the WorkOS
// session tokens in a cookie the client cannot read or tamper with.
package securecookie

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
)

// ErrInvalid is returned when a sealed value cannot be decrypted (wrong key,
// truncated, or tampered with).
var ErrInvalid = errors.New("securecookie: invalid or tampered value")

// Seal encrypts plaintext with key (must be 32 bytes) and returns a
// URL-safe base64 string of nonce||ciphertext.
func Seal(key, plaintext []byte) (string, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.RawURLEncoding.EncodeToString(sealed), nil
}

// Unseal reverses Seal. It returns ErrInvalid if the value does not
// authenticate against key.
func Unseal(key []byte, value string) ([]byte, error) {
	gcm, err := newGCM(key)
	if err != nil {
		return nil, err
	}
	raw, err := base64.RawURLEncoding.DecodeString(value)
	if err != nil {
		return nil, ErrInvalid
	}
	if len(raw) < gcm.NonceSize() {
		return nil, ErrInvalid
	}
	nonce, ciphertext := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, ErrInvalid
	}
	return plaintext, nil
}

func newGCM(key []byte) (cipher.AEAD, error) {
	if len(key) == 0 {
		return nil, errors.New("securecookie: key must not be empty")
	}
	// Derive a fixed 32-byte AES-256 key from the configured password so any
	// non-empty password works (no exact-length requirement).
	sum := sha256.Sum256(key)
	block, err := aes.NewCipher(sum[:])
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}
