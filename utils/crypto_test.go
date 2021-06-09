package utils

import (
	"crypto/aes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCBCCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	iv := key[:aes.BlockSize]
	plainText := "Iloveyiigo"

	// ZERO_PADDING
	zero := NewCBCCrypto(key, iv, ZERO)

	e0b, err := zero.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d0b, err := zero.Decrypt(e0b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d0b))

	// PKCS5_PADDING
	pkcs5 := NewCBCCrypto(key, iv, PKCS5)

	e5b, err := pkcs5.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d5b, err := pkcs5.Decrypt(e5b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d5b))

	// PKCS7_PADDING
	pkcs7 := NewCBCCrypto(key, iv, PKCS7)

	e7b, err := pkcs7.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d7b, err := pkcs7.Decrypt(e7b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d7b))
}

func TestECBCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	plainText := "Iloveyiigo"

	// ZERO_PADDING
	zero := NewECBCrypto(key, ZERO)

	e0b, err := zero.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d0b, err := zero.Decrypt(e0b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d0b))

	// PKCS5_PADDING
	pkcs5 := NewECBCrypto(key, PKCS5)

	e5b, err := pkcs5.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d5b, err := pkcs5.Decrypt(e5b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d5b))

	// PKCS7_PADDING
	pkcs7 := NewECBCrypto(key, PKCS7)

	e7b, err := pkcs7.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	d7b, err := pkcs7.Decrypt(e7b)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(d7b))
}

func TestCFBCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	iv := key[:aes.BlockSize]
	plainText := "Iloveyiigo"

	cfb := NewCFBCrypto(key, iv)

	eb, err := cfb.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	db, err := cfb.Decrypt(eb)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(db))
}

func TestOFBCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	iv := key[:aes.BlockSize]
	plainText := "Iloveyiigo"

	ofb := NewOFBCrypto(key, iv)

	eb, err := ofb.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	db, err := ofb.Decrypt(eb)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(db))
}

func TestCTRCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	iv := key[:aes.BlockSize]
	plainText := "Iloveyiigo"

	ctr := NewCTRCrypto(key, iv)

	eb, err := ctr.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	db, err := ctr.Decrypt(eb)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(db))
}

func TestGCMCrypto(t *testing.T) {
	key := []byte("AES256Key-32Characters1234567890")
	nonce := key[:12]
	plainText := "Iloveyiigo"

	gcm := NewGCMCrypto(key, nonce)

	eb, err := gcm.Encrypt([]byte(plainText))
	assert.Nil(t, err)

	db, err := gcm.Decrypt(eb)
	assert.Nil(t, err)
	assert.Equal(t, plainText, string(db))
}

//func TestRSASign(t *testing.T) {
//	plainText := "Iloveyiigo"
//
//	signature, err := RSASignWithSha256([]byte(plainText), privateKey)
//
//	assert.Nil(t, err)
//	assert.Nil(t, RSAVerifyWithSha256([]byte(plainText), signature, publicKey))
//}
//
//func TestRSACrypto(t *testing.T) {
//	plainText := "Iloveyiigo"
//
//	eb, err := RSAEncrypt([]byte(plainText), publicKey)
//
//	assert.Nil(t, err)
//
//	db, err := RSADecrypt(eb, privateKey)
//
//	assert.Nil(t, err)
//	assert.Equal(t, plainText, string(db))
//
//	eboeap, err := RSAEncryptOEAP([]byte(plainText), publicKey)
//
//	assert.Nil(t, err)
//
//	dboeap, err := RSADecryptOEAP(eboeap, privateKey)
//
//	assert.Nil(t, err)
//	assert.Equal(t, plainText, string(dboeap))
//}
