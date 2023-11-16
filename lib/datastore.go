package lib

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
)

type DataStore struct {
	Users              []User // List of users
	Files              []File // List of Files
	UseBackdoorKey     bool   // Flag for Backdoor
	BackdoorPubKeyPath string //Path to Backdoor PublicKey
}

type User struct {
	Id                   string // User UUID
	UserName             string // Username
	PublicKeyPath        string // Path to Public Key
	PrivateKeyPath       string // Path to Private Key
	BackdoorKeySignature string // Signature of the Backdoor key
	//TrustedKeysSignatures    []TrustedKeySignature // List of Signed Public Keys
	TrustedKeysSignaturesMap map[string]TrustedKeySignature

	// SignatureStoreSignature string           // Signature of Hash(BackdoorKey + Trusted Keys)
}

type SharedUserKey struct {
	UserName     string // Username
	EncryptedKey string // Encrypted AES Key
}

type TrustedKeySignature struct {
	UserName  string // Username
	Signature string // Signature of UserName + PublicKey
}

type File struct {
	Id             string          // UUID of File
	Name           string          // Name of File
	Owner          string          // Username of File Owner
	FilePath       string          // File Path of Encrypted File
	EncryptedKey   string          // Encrypted AES Key
	SharedUserKeys []SharedUserKey // Username + Encrypted Key for Shared User
}

func (dataStore *DataStore) WriteDataStoreToFile(path string) {

	j, err := json.Marshal(*dataStore)
	CheckError(err)

	err = os.WriteFile(path+"/vault.json", j, 0644)
	CheckError(err)

}

func (dataStore *DataStore) LoadDataStoreFromFile(path string) {

	content, err := os.ReadFile(path + "/vault.json")
	checkError(err)

	err = json.Unmarshal(content, dataStore)
	checkError(err)

}

func (dataStore *DataStore) AddUser(userName string, dataPath string) bool {

	if userExists(dataStore, userName) {
		return false
	}

	var user User

	user.Id = uuid.NewString()
	user.UserName = userName
	user.PublicKeyPath = dataPath + "/keys/" + userName + "_" + user.Id + ".pub"
	user.PrivateKeyPath = dataPath + "/keys/" + userName + "_" + user.Id + ".priv"

	keyPair := GenerateRSAKeypair()

	StoreRSAKeypairDialog(user.PrivateKeyPath, user.PublicKeyPath, keyPair)

	user.TrustedKeysSignaturesMap = make(map[string]TrustedKeySignature)

	for i := 0; i < len(dataStore.Users); i++ {
		signature := dataStore.Users[i].CreateUserPublicKeySignature(keyPair)
		//user.TrustedKeysSignatures = append(user.TrustedKeysSignatures, signature)
		user.TrustedKeysSignaturesMap[dataStore.Users[i].UserName] = signature
	}

	user.BackdoorKeySignature = dataStore.createBackdoorSignature(keyPair)

	// user.SignatureStoreSignature = user.CreateSignatureStoreSignature(keyPair)

	dataStore.Users = append(dataStore.Users, user)

	return true

}

func (user *User) CreateUserPublicKeySignature(keyPair *rsa.PrivateKey) TrustedKeySignature {
	var signature TrustedKeySignature

	content, err := os.ReadFile(user.PublicKeyPath)

	if err != nil {
		fmt.Println("Failed To Read Private Key from User")
		fmt.Println(err)
		os.Exit(2)
	}

	signature.UserName = user.UserName
	signatureString := user.GeneratePublicKeySignatureString(content)
	signature.Signature = base64.StdEncoding.EncodeToString(SignWithHash(keyPair, signatureString))

	return signature
}

func (user *User) GeneratePublicKeySignatureString(content []byte) string {

	return "Key::" + user.UserName + "::" + base64.StdEncoding.EncodeToString(content)

}

func (dataStore *DataStore) createBackdoorSignature(keyPair *rsa.PrivateKey) string {

	content, err := os.ReadFile(dataStore.BackdoorPubKeyPath)

	if err != nil {
		fmt.Println("Failed To Read Backdoor Key")
		fmt.Println(err)
		os.Exit(2)
	}

	hash := base64.StdEncoding.EncodeToString(SignWithHash(keyPair, "Backdoor-Key::"+string(content)))

	return hash
}

func (dataStore *DataStore) GetUserArrId(userName string) int {

	for i := 0; i < len(dataStore.Users); i++ {

		if userName == dataStore.Users[i].UserName {
			return i
		}

	}

	return -1
}

func (dataStore *DataStore) DeleteUser(userName string) {

	var _Users []User

	for i := 0; i < len(dataStore.Users); i++ {

		if userName != dataStore.Users[i].UserName {
			_Users = append(_Users, dataStore.Users[i])
		}

	}

	dataStore.Users = _Users

}

func userExists(dataStore *DataStore, userName string) bool {

	for i := 0; i < len(dataStore.Users); i++ {

		if userName == dataStore.Users[i].UserName {
			return true
		}

	}

	return false
}

func (user *User) PrintUser() {

	if user == nil || &user == nil {
		fmt.Println("User to print is nil")
		os.Exit(2)
	}
	fmt.Printf("Id: %s\nUsername: %s\n", user.Id, user.UserName)

}

// Deprecated
func (user *User) CreateSignatureStoreSignature(keyPair *rsa.PrivateKey) string {

	_user := *user
	//_user.SignatureStoreSignature = "IGNORE"

	j, err := json.Marshal(_user)
	if err != nil {
		fmt.Println("Failed to generate Json for Hash")
		fmt.Println(err)
		os.Exit(2)
	}

	hash := base64.StdEncoding.EncodeToString(SignWithHash(keyPair, string(j)))

	return hash
}
