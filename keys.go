package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
)

var (
	pubKeyFilePath    = filepath.Join(os.Getenv("HOME"), ".fir", "keys", "fir-publickey.pub")
	privKeyFilePath   = filepath.Join(os.Getenv("HOME"), ".fir", "keys", "fir-privatekey.private")
	signedKeyFilePath = filepath.Join(os.Getenv("HOME"), ".fir", "keys", "fir-signedkey.sig")
	selfCertFilePath  = filepath.Join(os.Getenv("HOME"), ".fir", "keys", "fir-signedcert.sig")
)

// ED25519Keys This is a struct for holding keys and a signature.
type ED25519Keys struct {
	publicKey  string
	privateKey string
	signedKey  string
	selfCert   string
}

func handle(msg string, err error) {
	if err != nil {
		fmt.Printf("\n%s: %s", msg, err)
	}
}

func fileExistsKey(filename string) bool {
	referencedFile, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !referencedFile.IsDir()
}

func readFileKey(filename string) string {
	text, err := os.ReadFile(filename)
	handle("Couldnt read the file: ", err)
	return string(text)
}

func createFileKey(filename string) {
	var _, err = os.Stat(filename)
	if os.IsNotExist(err) {
		var file, err = os.Create(filename)
		handle("", err)
		defer file.Close()
	}
}

func writeFileKey(filename, textToWrite string) {
	var file, err = os.OpenFile(filename, os.O_RDWR, 0644)
	handle("", err)
	defer file.Close()
	_, err = file.WriteString(textToWrite)
	err = file.Sync()
	handle("", err)
}

func initKeys() *ED25519Keys {
	if !fileExistsKey(privKeyFilePath) {
		generateKeys()
	}
	keys := ED25519Keys{}
	keyspublicKey := readFileKey(pubKeyFilePath)
	keysprivateKey := readFileKey(privKeyFilePath)
	keyssignedKey := readFileKey(signedKeyFilePath)
	keysselfCert := readFileKey(selfCertFilePath)
	keys.publicKey = keyspublicKey[:64]
	keys.privateKey = keysprivateKey[:64]
	keys.signedKey = keyssignedKey[:64]
	keys.selfCert = keysselfCert[:64]
	return &keys
}

func generateKeys() *ED25519Keys {
	keys := ED25519Keys{}
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		handle("error: ", err)
	}
	keys.privateKey = hex.EncodeToString(privKey[0:32])
	keys.publicKey = hex.EncodeToString(pubKey)
	signedKey := ed25519.Sign(privKey, pubKey)
	keys.signedKey = hex.EncodeToString(signedKey)
	keys.selfCert = keys.publicKey + keys.signedKey
	createFileKey(pubKeyFilePath)
	createFileKey(privKeyFilePath)
	createFileKey(signedKeyFilePath)
	createFileKey(selfCertFilePath)
	writeFileKey(pubKeyFilePath, keys.publicKey[:64])
	writeFileKey(privKeyFilePath, keys.privateKey[:64])
	writeFileKey(signedKeyFilePath, keys.signedKey[:64])
	writeFileKey(selfCertFilePath, keys.selfCert[:64])
	return &keys
}

func sign(myKeys *ED25519Keys, msg string) string {
	messageBytes := []byte(msg)
	privateKey, err := hex.DecodeString(myKeys.privateKey)
	if err != nil {
		handle("private key error: ", err)
	}
	publicKey, err := hex.DecodeString(myKeys.publicKey)
	if err != nil {
		handle("public key error: ", err)
	}
	privateKey = append(privateKey, publicKey...)
	signature := ed25519.Sign(privateKey, messageBytes)
	return hex.EncodeToString(signature)
}

func signKey(myKeys *ED25519Keys, publicKey string) string {
	messageBytes, err := hex.DecodeString(publicKey)
	if err != nil {
		handle("error: ", err)
	}
	privateKey, err := hex.DecodeString(myKeys.privateKey)
	if err != nil {
		handle("error: ", err)
	}
	pubKey, err := hex.DecodeString(myKeys.publicKey)
	if err != nil {
		handle("error: ", err)
	}
	privateKey = append(privateKey, pubKey...)
	signature := ed25519.Sign(privateKey, messageBytes)
	return hex.EncodeToString(signature)
}

func verifySignature(publicKey string, msg string, signature string) bool {
	pubKey, err := hex.DecodeString(publicKey)
	if err != nil {
		handle("error: ", err)
	}
	messageBytes := []byte(msg)
	sig, err := hex.DecodeString(signature)
	if err != nil {
		handle("error: ", err)
	}
	return ed25519.Verify(pubKey, messageBytes, sig)
}

func verifySignedKey(publicKey string, publicSigningKey string, signature string) bool {
	pubKey, err := hex.DecodeString(publicKey)
	if err != nil {
		handle("error: ", err)
	}
	pubSignKey, err := hex.DecodeString(publicSigningKey)
	if err != nil {
		handle("error: ", err)
	}
	sig, err := hex.DecodeString(signature)
	if err != nil {
		handle("error: ", err)
	}
	return ed25519.Verify(pubSignKey, pubKey, sig)
}
