package rsautil

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
)

// SignSha256 sha256签名
func SignSha256(data []byte, key *rsa.PrivateKey) ([]byte, error) {
	return Sign(data, key, crypto.SHA256)
}

func VerifySha256(plainData []byte, signature []byte, publicKey *rsa.PublicKey) error {
	return Verify(plainData, signature, publicKey, crypto.SHA256)
}

// Verify hashed data with rsa
func Verify(data []byte, signature []byte, pub *rsa.PublicKey, hash crypto.Hash) error {
	var h = hash.New()
	h.Write(data)
	return rsa.VerifyPKCS1v15(pub, hash, h.Sum(nil), signature)
}

// Sign rsa 签名
func Sign(data []byte, key *rsa.PrivateKey, hash crypto.Hash) ([]byte, error) {
	var h = hash.New()
	h.Write(data)
	var hashed = h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, key, hash, hashed)
}
