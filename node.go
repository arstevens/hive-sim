package main

import (
  "log"
  "crypto"
  "crypto/rsa"
  "crypto/rand"
  "crypto/sha256"
)

type Node interface {
  Encrypt([]byte) []byte
  Decrypt([]byte) []byte
  Sign([]byte) []byte
  Verify([]byte) []byte
}

type BaseNode struct {
  secretKey rsa.PrivateKey
}

func (bn BaseNode) Encrypt(msg []byte) []byte {
  pubKey := bn.secretKey.Public()
  rsaPubKey := pubKey.(*rsa.PublicKey)

  crypt, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, rsaPubKey, msg, []byte{})
  if err != nil {
    log.Print(err)
    return []byte{}
  }
  return crypt
}

func (bn BaseNode) Decrypt(cipherTxt []byte) []byte {
  privKey := bn.secretKey

  msg, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, &privKey, cipherTxt, []byte{})
  if err != nil {
    log.Print(err)
    return []byte{}
  }
  return msg
}

func (bn BaseNode) Sign(data []byte) []byte {
  privKey := bn.secretKey

  hashedArray := sha256.Sum256(data)
  hashedSlice := hashedArray[:]

  sign, err := rsa.SignPKCS1v15(rand.Reader, &privKey, crypto.SHA256, hashedSlice)
  if err != nil {
    log.Print(err)
    return []byte{}
  }
  return sign
}

func (bn BaseNode) Verify(data []byte, sig []byte) bool {
  pubKey := bn.secretKey.Public()
  rsaPubKey := pubKey.(*rsa.PublicKey)

  hashedArray := sha256.Sum256(data)
  hashedSlice := hashedArray[:]

  err := rsa.VerifyPKCS1v15(rsaPubKey, crypto.SHA256, hashedSlice, sig)
  if err != nil {
    log.Print(err)
    return false
  }
  return true
}
