package lib

import (
	"crypto"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"golang.org/x/crypto/argon2"
	"io"
	"log"
	"os"
)

func GenerateAESKey(key *[]byte) {
	_, err := rand.Reader.Read(*key)
	checkError(err)
}

func GenerateAESKeyPointer() *[]byte {

	key := make([]byte, keyLengthAES)

	GenerateAESKey(&key)

	return &key

}

func CleanMemory(mem *[]byte) {
	for i := range *mem {
		(*mem)[i] = 0
	}
}

func EncryptAESGCM(key *[]byte, plaintext []byte) []byte {

	c, err := aes.NewCipher(*key)
	checkError(err)

	gcm, err := cipher.NewGCM(c) // use AES-GCM
	checkError(err)

	// Generate random Nonce (Number used once)
	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	checkError(err)

	// prepended Nonce/IV || Ciphertext (+Tag)
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext
}

func DecryptAESGCM(key *[]byte, ciphertext []byte) []byte {

	c, err := aes.NewCipher(*key)
	checkError(err)

	log.Println("Going to GCM")

	gcm, err := cipher.NewGCM(c) // use AES-GCM
	checkError(err)

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		log.Fatal("ciphertext < nonce")
	}
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		log.Println("Failed to Decrypt.")
		log.Fatal(err)
	}

	return plaintext
}

func GenA2IDKey(passwd *[]byte, salt []byte) *[]byte {

	key := make([]byte, keyLengthPasswd)

	key = argon2.IDKey(*passwd, salt, 1, 64*1024, 4, keyLengthPasswd)

	return &key

}

func EncryptFile(key *[]byte, fPath string, newPath string) {

	content, err := os.ReadFile(fPath)
	checkError(err)

	ciphertext := EncryptAESGCM(key, content)

	err = os.WriteFile(newPath, ciphertext, 0640)
	checkError(err)
}

func DecryptFile(key *[]byte, fPath string, newPath string) {

	content, err := os.ReadFile(fPath)
	checkError(err)

	cleartext := DecryptAESGCM(key, content)

	err = os.WriteFile(newPath, cleartext, 0640)
	checkError(err)
}

func EncryptFileWithBackdoor(key *[]byte, fPath string, newPath string, masterKey *rsa.PrivateKey) {

	// TODO
}

func DecryptFileWithBackdoor(key *[]byte, fPath string, newPath string) {

	// TODO
}

func DecryptFileBackdoor(masterKey *rsa.PrivateKey, fPath string, newPath string) {

	// TODO
}

func GenerateRSAKeypair() *rsa.PrivateKey {

	keyPair, err := rsa.GenerateKey(rand.Reader, 4096)
	checkError(err)

	// validate key
	err = keyPair.Validate()
	if err != nil {
		log.Printf("ERROR: fail to validate key pair, %s", err.Error())
		os.Exit(1)
	}

	return keyPair

}

func StoreRSAKeypair(fPathPrivateKey string, fPathPubKey string, keyPair *rsa.PrivateKey, keyPasswd *[]byte, salt []byte) {

	privateKey := x509.MarshalPKCS1PrivateKey(keyPair)
	//publicKey := x509.MarshalPKCS1PublicKey(&keyPair.PublicKey)

	encryptedPrivateKey := EncryptAESGCM(keyPasswd, privateKey)

	encryptedPrivateKey = append(salt, encryptedPrivateKey...)

	pemPrivateKey := []byte(base64.StdEncoding.EncodeToString(encryptedPrivateKey))

	err := os.WriteFile(fPathPrivateKey, pemPrivateKey, 0644)
	checkError(err)

	StoreRSAPublicKey(fPathPubKey, &(*keyPair).PublicKey)

	//err = os.WriteFile(fPath+".pub", []byte(base64.StdEncoding.EncodeToString(publicKey)), 0644)
	//checkError(err)

}

func StoreRSAPublicKey(fPath string, key *rsa.PublicKey) {

	publicKey := x509.MarshalPKCS1PublicKey(key)

	err := os.WriteFile(fPath, []byte(base64.StdEncoding.EncodeToString(publicKey)), 0644)
	checkError(err)
}

func LoadRSAKeypair(fPath string) *rsa.PrivateKey {

	content, err := os.ReadFile(fPath)
	checkError(err)

	pemPrivateKey, err := base64.StdEncoding.DecodeString(string(content))
	checkError(err)

	salt := pemPrivateKey[:saltLength]
	pemPrivateKey = pemPrivateKey[saltLength:]

	fmt.Print("Please Enter your Password: ")
	keyPasswd := ReadPasswd()

	key := GenA2IDKey(keyPasswd, salt)
	CleanMemory(keyPasswd)

	log.Println("Generated Argon2")

	decryptedPEM := DecryptAESGCM(key, pemPrivateKey)

	log.Println("Decrypted Key")

	privateKey, err := x509.ParsePKCS1PrivateKey(decryptedPEM)
	checkError(err)

	err = privateKey.Validate()
	if err != nil {
		log.Fatalf("ERROR: fail to validate loaded key pair, %s", err.Error())
	}

	return privateKey
}

func LoadRSAPublicKey(fPath string) *rsa.PublicKey {
	content, err := os.ReadFile(fPath)
	checkError(err)

	_publicKey, err := base64.StdEncoding.DecodeString(string(content))
	checkError(err)

	publicKey, err := x509.ParsePKCS1PublicKey(_publicKey)
	checkError(err)

	return publicKey
}

func SignWithHash(keyPair *rsa.PrivateKey, message string) []byte {
	digest := HashString(message)

	signature, err := rsa.SignPSS(rand.Reader, keyPair, crypto.SHA512, digest, nil)
	checkError(err)

	return signature

}

func SignatureIsValid(pubKey *rsa.PublicKey, message string, signature []byte) bool {
	digest := HashString(message)

	err := rsa.VerifyPSS(pubKey, crypto.SHA512, digest, signature, nil)

	if err != nil {
		log.Println("Failed to Verify Signature")
		return false
	}
	log.Println("Verified Signature")
	return true
}

func EncryptWithRSA(pubKey *rsa.PublicKey, plaintext *[]byte) []byte {

	ciphertext, err := rsa.EncryptOAEP(sha512.New(), rand.Reader, pubKey, *plaintext, nil)
	checkError(err)
	return ciphertext
}

func DecryptWithRSA(keyPair *rsa.PrivateKey, ciphertext []byte) *[]byte {

	plaintext, err := keyPair.Decrypt(nil, ciphertext, &rsa.OAEPOptions{Hash: crypto.SHA512})
	checkError(err)
	return &plaintext
}

func HashString(str string) []byte {
	h := sha512.New()
	h.Write([]byte(str))
	return h.Sum(nil)
}
