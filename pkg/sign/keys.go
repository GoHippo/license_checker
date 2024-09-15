package sign

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"fmt"
)

// Функция для создания пары ключей (приватного и публичного)
func generateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// Преобразование публичного ключа в формат DER
func publicKeyToHex(pubKey *ecdsa.PublicKey) (string, error) {
	key, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", fmt.Errorf("marshal public key to DER err:%v", err)
	}

	return hex.EncodeToString(key), nil
}

func parseHexPublicKey(hexKey string) (*ecdsa.PublicKey, error) {

	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}

	pubKey, err := x509.ParsePKIXPublicKey(key)
	if err != nil {
		return nil, err
	}

	switch pubKey := pubKey.(type) {
	case *ecdsa.PublicKey:
		return pubKey, nil
	default:
		return nil, fmt.Errorf("неизвестный тип ключа: %T", pubKey)
	}
}
