package ecdsa

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"errors"
)

// DefaultCurve は鍵生成時に使用する既定の曲線です。
var DefaultCurve = elliptic.P256()

// GenerateKey は楕円曲線と暗号学的乱数生成器から秘密鍵を生成します。
func GenerateKey(curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	if curve == nil {
		curve = DefaultCurve
	}
	return ecdsa.GenerateKey(curve, rand.Reader)
}

// Sign は秘密鍵でメッセージを署名し ASN.1 形式のシグネチャを返します。
// メッセージは SHA-256 でハッシュ化されます。
func Sign(priv *ecdsa.PrivateKey, message []byte) ([]byte, error) {
	if priv == nil {
		return nil, errors.New("ecdsa: private key is nil")
	}
	hash := sha256.Sum256(message)
	return ecdsa.SignASN1(rand.Reader, priv, hash[:])
}

// Verify は公開鍵で署名を検証します。メッセージは Sign と同様にハッシュ化されます。
func Verify(pub *ecdsa.PublicKey, message, signature []byte) bool {
	if pub == nil {
		return false
	}
	hash := sha256.Sum256(message)
	return ecdsa.VerifyASN1(pub, hash[:], signature)
}

// PublicKey は秘密鍵から公開鍵を取り出します。
func PublicKey(priv *ecdsa.PrivateKey) (*ecdsa.PublicKey, error) {
	if priv == nil {
		return nil, errors.New("ecdsa: private key is nil")
	}
	return &priv.PublicKey, nil
}
