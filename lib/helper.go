package lib

import (
	"crypto/rsa"
	"crypto/sha512"
	"fmt"
	"golang.org/x/crypto/ssh/terminal"
	"log"
	"os"
	"path"
)

func getHash(x []byte) []byte {
	h := sha512.New()
	_, err := h.Write(x)
	if err != nil {
		log.Println("Failed to Hash")
		log.Fatal(err)
	}
	return h.Sum(nil)
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func CheckError(err error) {
	checkError(err)
}

func ReadPasswd() *[]byte {
	passwd, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	fmt.Print("\n")
	checkError(err)
	return &passwd
}

func StoreRSAKeypairDialog(fPathPrivateKey string, fPathPubKey string, keyPair *rsa.PrivateKey) {

	salt := make([]byte, saltLength)
	GenerateAESKey(&salt)

	fmt.Print("Please Enter your Password: ")
	passwd := ReadPasswd()

	passwdKey := GenA2IDKey(passwd, salt)
	CleanMemory(passwd)

	StoreRSAKeypair(fPathPrivateKey, fPathPubKey, keyPair, passwdKey, salt)

	CleanMemory(passwdKey)

}

func LoadRSAKeypairDialog(fPathPrivateKey string) *rsa.PrivateKey {

	return LoadRSAKeypair(fPathPrivateKey)

}

func BuildEncryptedFileName(fPath string, fileUUID string) string {
	x := fileUUID + "_" + path.Base(fPath) + ".encrypted"
	return x
}
