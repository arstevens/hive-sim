package controller

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"log"
)

type RsaNode struct {
	secretKey *rsa.PrivateKey
}

func (rn RsaNode) Encrypt(msg []byte) []byte {
	pubKey := rn.secretKey.Public()
	rsaPubKey := pubKey.(*rsa.PublicKey)

	crypt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, msg, []byte{})
	if err != nil {
		log.Print(err)
		return []byte{}
	}
	return crypt
}

func (rn RsaNode) Decrypt(cipherTxt []byte) []byte {
	privKey := rn.secretKey

	msg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privKey, cipherTxt, []byte{})
	if err != nil {
		log.Print(err)
		return []byte{}
	}
	return msg
}

func (rn RsaNode) Sign(data []byte) []byte {
	privKey := rn.secretKey

	hashedArray := sha256.Sum256(data)
	hashedSlice := hashedArray[:]

	sign, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashedSlice)
	if err != nil {
		log.Print(err)
		return []byte{}
	}
	return sign
}

func (rn RsaNode) PublicKey() *rsa.PublicKey {
	return rn.secretKey.Public().(*rsa.PublicKey)
}

func RsaVerify(data []byte, sig []byte, pubKey *rsa.PublicKey) bool {
	hashedArray := sha256.Sum256(data)
	hashedSlice := hashedArray[:]

	err := rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashedSlice, sig)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}
