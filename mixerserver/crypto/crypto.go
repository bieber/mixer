/*
 * Copyright 2015, Robert Bieber
 *
 * This file is part of mixer.
 *
 * mixer is free software: you can redistribute it and/or modify it
 * under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * mixer is distributed in the hope that it will be useful,
 * WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with mixer.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package crypto implements AES encryption for CSRF tokens and API
// secrets to send down to the client.  That way we can avoid having
// to maintain any server-side authentication state without just
// handing the OAuth tokens to the users.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
)

var aesKey []byte

// GenerateAESKey generates a new random AES key, and returns it
// encoded as base64.
func GenerateAESKey() (string, error) {
	keyBytes := make([]byte, 16, 16)
	_, err := rand.Read(keyBytes)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(keyBytes), nil
}

// SetAESKey sets the AES key to use for crypto operations.  It should
// be 16 bytes encoded in base64.
func SetAESKey(key string) error {
	var err error
	aesKey, err = base64.URLEncoding.DecodeString(key)
	if err != nil {
		return err
	}

	if len(aesKey) != 16 {
		return errors.New("Invalid key length")
	}
	return nil
}

// Encrypt encrypts a string and returns the ciphertext as a base64
// encoded string.  SetAESKey must have been called previously.
func Encrypt(plaintext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	toEncrypt := []byte(plaintext)
	ciphertext := make([]byte, len(toEncrypt)+aes.BlockSize)

	iv := ciphertext[:aes.BlockSize]
	_, err = rand.Read(iv)
	if err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], toEncrypt)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext (encoded as base64) and attempts to
// return the decoded value as a string.  SetAESKey must have been
// called previously.
func Decrypt(ciphertext string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	toDecrypt, err := base64.URLEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	if len(toDecrypt) < aes.BlockSize {
		return "", errors.New("Ciphertext is too short")
	}

	iv := toDecrypt[:aes.BlockSize]
	plaintext := make([]byte, len(toDecrypt)-aes.BlockSize)

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(plaintext, toDecrypt[aes.BlockSize:])

	return string(plaintext), nil
}
