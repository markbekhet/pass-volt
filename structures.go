package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
)

type AccountDetails struct {
	Username string
	Password []byte
}

type LoginInfo struct {
	Id      string
	Details AccountDetails
}

type Accounts struct {
	Accounts []LoginInfo
}

type Encrypter interface {
	encrypt()
}

type Decrypter interface {
	decrypt()
}

func (d *AccountDetails) encrypt() {
	encrypt := encryptHelper(d.Password, d.Username)
	d.Password = encrypt

}

func (d AccountDetails) decrypt() AccountDetails {
	var ret AccountDetails
	ret.Username = d.Username
	ret.Password = decryptHelper(d.Password, d.Username)
	return ret
}

func shaHashing(input string) string {
	plainText := []byte(input)
	hash := md5.Sum(plainText)
	return hex.EncodeToString(hash[:])
}

func encryptHelper(value []byte, keyPhrase string) []byte {
	aesBlock, err := aes.NewCipher([]byte(shaHashing(keyPhrase)))
	if err != nil {
		fmt.Println(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		fmt.Println(err)
	}
	nonce := make([]byte, gcmInstance.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipheredText := gcmInstance.Seal(nonce, nonce, []byte(value), nil)

	return cipheredText
}

func decryptHelper(ciphered []byte, keyphrase string) []byte {
	hash := shaHashing(keyphrase)
	aesBlock, err := aes.NewCipher([]byte(hash))
	if err != nil {
		log.Fatalln(err)
	}
	gcmInstance, err := cipher.NewGCM(aesBlock)
	if err != nil {
		log.Fatalln(err)
	}
	nonceSize := gcmInstance.NonceSize()
	nonce, cipheredText := ciphered[:nonceSize], ciphered[nonceSize:]
	originalText, err := gcmInstance.Open(nil, nonce, cipheredText, nil)
	if err != nil {
		log.Fatalln(err)
	}
	return originalText
}
