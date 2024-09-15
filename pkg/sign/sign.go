package sign

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"math/big"
)

// Функция для создания цифровой подписи
func createSign(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	hashed := sha256.Sum256([]byte(message))
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, hashed[:])

	sign := append(r.Bytes(), s.Bytes()...)
	signStr := base64.StdEncoding.EncodeToString(sign)
	return signStr, err
}

// Функция для проверки цифровой подписи
func VerifySign(message string, sign string, hexPublicKey string) (bool, error) {
	var op = `sign.VerifySign`

	signatureBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return false, fmt.Errorf("%s: decode signature: %w", op, err)
	}

	publicKey, err := parseHexPublicKey(hexPublicKey)
	if err != nil {
		return false, fmt.Errorf("%s: parse public key: %w", op, err)
	}

	half := len(signatureBytes) / 2
	r := new(big.Int).SetBytes(signatureBytes[:half])
	s := new(big.Int).SetBytes(signatureBytes[half:])

	hashed := sha256.Sum256([]byte(message))

	return ecdsa.Verify(publicKey, hashed[:], r, s), nil
}

func NewCreateSign(data string) (string, string, error) {
	var op = `pkg.sign.NewCreateSign`

	// Создание пары ключей
	privateKey, publicKey, err := generateKeyPair()
	if err != nil {
		return "", "", fmt.Errorf("%s: Create keys: %w", op, err)
	}

	// Подписание сообщения
	sign, err := createSign(data, privateKey)
	if err != nil {
		return "", "", fmt.Errorf("%s: Create Sign: %w", op, err)
	}

	hexPublicKey, err := publicKeyToHex(publicKey)
	if err != nil {
		return "", "", fmt.Errorf("%s:%w", op, err)
	}

	return sign, hexPublicKey, nil
}
